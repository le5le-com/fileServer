package file

import (
	"errors"
	"fileServer/utils"
	"os"
	"path"
	"strings"

	"fileServer/keys"

	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	if !IsExist(CachePath) {
		err := os.MkdirAll(CachePath, os.ModePerm)
		if err != nil {
			log.Panic().Err(err).Msgf("Fail to create the CachePath: %s", CachePath)
		}
	}
}

// Limit 普通文件大小限制
func Limit(ctx iris.Context) {
	if ctx.GetContentLength() > 10<<20 {
		ctx.StatusCode(iris.StatusRequestEntityTooLarge)
		ctx.JSON(bson.M{"error": keys.ErrorFileMaxSize})
		return
	}
	ctx.Next()
}

// UploadFile 上传普通文件
func UploadFile(ctx iris.Context) {
	path := upload(ctx, false)
	if path != "" {
		ctx.JSON(bson.M{"url": "/file" + path})
	}
}

// DelFile 删除文件
func DelFile(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	filename := "/" + ctx.Params().Get("path")
	log.Debug().Msg(filename)
	if filename == "" {
		ret["error"] = errors.New(keys.ErrorParam)
		return
	}

	uid := ctx.Values().GetString("uid")
	err := Del(filename, CachePath+"/"+uid, uid)
	if err != nil {
		ret["error"] =  err.Error()
		return
	}

}

// UploadImage 上传图片
func UploadImage(ctx iris.Context) {
	path := upload(ctx, true)
	if path != "" {
		ctx.JSON(bson.M{"url": "/image" + path})
	}
}

// File 获取普通文件
func File(ctx iris.Context) {
	fullpath := download(ctx)
	if fullpath != "" {
		ctx.ServeFile(fullpath, true)
	}
}

// Image 获取图片
func Image(ctx iris.Context) {
	fullpath := download(ctx)
	if fullpath == "" {
		return
	}

	w, _ := ctx.URLParamInt("w")
	h, _ := ctx.URLParamInt("h")
	if w > 0 || h > 0 {
		thumb, err := ImageThumbnail(fullpath, w, h)
		if err == nil {
			fullpath = thumb
		} else {
			log.Error().Err(err).Str("func", "file.Image").Msgf("Fail to thumbnail: fullpath=%s", fullpath)
		}
	}

	ctx.ServeFile(fullpath, true)
}

// Images 获取目录图片
func Images(ctx iris.Context) {
	ctx.JSON(GetFileList(CachePath + "/materials/"))
}

// FileAttr 修改文件属性
func FileAttr(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	filename := "/" + ctx.Params().Get("path")
	if filename == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.Text("")
		return
	}

	data := bson.M{}
	err := ctx.ReadJSON(&data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		if err != nil {
			ret["errorDetail"] = err.Error()
		}
		return
	}

	uid := ctx.Values().GetString("uid")
	data["userId"] = uid
	data["username"] = ctx.Values().GetString("username")
	err = PatchAttr(filename, data, uid)
	if err != nil {
		ret["error"] = keys.ErrorSave
		if err != nil {
			ret["errorDetail"] = err.Error()
		}
	}
}

// upload 上传文件
func upload(ctx iris.Context, public bool) string {
	file, info, err := ctx.FormFile("file")
	if err != nil {
		log.Warn().Err(err).Str("func", "file.upload").Msg("Fail to read file from FormFile.")
		ctx.JSON(bson.M{"error": keys.ErrorFile, "system": err.Error()})
		return ""
	}
	defer file.Close()

	fullname := ctx.FormValue("path")
	if fullname == "" {
		fullname = info.Filename
	}
	if fullname[0] != '/' {
		fullname = "/" + fullname
	}

	randomName := ctx.FormValue("randomName")
	if randomName != "" {
		dir, filename := path.Split(fullname)
		ext := path.Ext(filename)
		rand, _ := utils.GetRandString(8)
		if ext != "" {
			fullname = dir + strings.TrimSuffix(filename, ext) + "_" + rand + ext
		} else {
			fullname = dir + filename + "_" + rand
		}
	}

	pub := ctx.FormValue("public")
	if pub == "false" {
		public = false
	} else if pub != "" {
		public = true
	}

	uid := ctx.Values().GetString("uid")
	fileInfo, _ := Info(fullname)
	if fileInfo != nil && fileInfo.Metadata.UserID == uid {
		ctx.JSON(bson.M{"error": keys.ErrorFileExists})
		return ""
	}

	tags := strings.Split(ctx.FormValue("tags"), ",")
	err = Put(fullname, file, bson.M{
		"userId":   ctx.Values().GetString("uid"),
		"username": ctx.Values().GetString("username"),
		"tags":     tags,
		"public":   public,
	})
	if err != nil {
		ctx.JSON(bson.M{"error": keys.ErrorFile, "system": err.Error()})
		return ""
	}

	return fullname
}

// download 获取文件
func download(ctx iris.Context) string {
	filename := "/" + ctx.Params().Get("path")
	if filename == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.ServeFile(CachePath+"/notFound.png", true)
		return ""
	}
	fullpath := CachePath + filename
	if IsExist(fullpath) {
		return fullpath
	}

	uid := ctx.Values().GetString("uid")
	fileInfo, _ := Info(filename)
	if fileInfo == nil {
		log.Warn().
			Str("func", "file.download").
			Str("remoteAddr", ctx.RemoteAddr()).
			Msgf("Error to read file: file=%s.", filename)
		ctx.StatusCode(iris.StatusNotFound)
		ctx.ServeFile(CachePath+"/notFound.png", true)
		return ""
	}
	if !fileInfo.Metadata.Public && fileInfo.Metadata.UserID != uid {
		log.Warn().
			Str("func", "file.download").
			Str("remoteAddr", ctx.RemoteAddr()).
			Msgf("No auth to read file: file=%s, fileinfo=%v.", filename, fileInfo)
		ctx.StatusCode(iris.StatusNotFound)
		ctx.ServeFile(CachePath+"/notFound.png", true)
		return ""
	}

	fullpath = CachePath + "/" + uid + filename
	if !IsExist(fullpath) {
		err := Get(filename, CachePath+"/"+uid)
		if err != nil {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.ServeFile(CachePath+"/notFound.png", true)
			return ""
		}
	}

	return fullpath
}
