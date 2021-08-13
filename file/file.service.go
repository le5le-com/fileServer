package file

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fileServer/config"
	"fileServer/db/mongo"
	"fileServer/keys"
)

const (
	// CachePath 文件临时存储目录
	CachePath = "./out/oss"
)

// WalkDir 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) ([]string, error) {
	files := make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err := filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		// 忽略目录
		if fi.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

// ReadFile 读取文件内容
func ReadFile(path string) (string, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()
	fileBytes, err := ioutil.ReadAll(fileHandle)
	return string(fileBytes), err
}

// IsExist 文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// GetUniqueName 获取一个唯一的文件名
func GetUniqueName() string {
	id := primitive.NewObjectID()
	newName := id.Hex()
	return newName
}

// Put 存储文件到数据库
func Put(name string, f interface{}, meta interface{}) error {
	if f == nil {
		log.Error().Str("func", "file.Put").Msg("文件指针为空.")
		return errors.New(keys.ErrorFileInfo)
	}

	bucketOptions := options.GridFSBucket().SetName(mongo.Files)
	bucket, _ := gridfs.NewBucket(mongo.Client.Database(config.App.Mongo.Database), bucketOptions)
	uploadOpts := options.GridFSUpload().SetMetadata(meta)
	uploadStream, err := bucket.OpenUploadStream(name, uploadOpts)
	if err != nil {
		log.Error().Err(err).Msg("Fail to openUploadStream on mongo.")
		return err
	}
	defer uploadStream.Close()

	switch data := f.(type) {
	case []byte:
		_, err = uploadStream.Write(data)
	case io.Reader:
		_, err = io.Copy(uploadStream, data)
	default:
	}

	if err != nil {
		log.Error().Err(err).Str("func", "file.Put").Msg("Fail to write file on mongo.")
	}

	return err
}

// Info 从数据库读取文件基础信息
func Info(name string) (*FileInfo, error) {
	if name == "" {
		return nil, errors.New(keys.ErrorParam)
	}

	fileInfo := &FileInfo{}
	bucketOptions := options.GridFSBucket().SetName(mongo.Files)
	bucket, _ := gridfs.NewBucket(mongo.Client.Database(config.App.Mongo.Database), bucketOptions)
	downloadStream, err := bucket.OpenDownloadStreamByName(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		downloadStream.Close()
	}()

	fileInfo.Filename = name
	err = bson.Unmarshal(downloadStream.GetFile().Metadata, &fileInfo.Metadata)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

// Get 从数据库读取文件
func Get(name, dir string) error {
	if name == "" {
		return errors.New(keys.ErrorParam)
	}
	bucketOptions := options.GridFSBucket().SetName(mongo.Files)
	bucket, _ := gridfs.NewBucket(mongo.Client.Database(config.App.Mongo.Database), bucketOptions)
	downloadStream, err := bucket.OpenDownloadStreamByName(name)
	if err != nil {
		log.Error().Err(err).Str("func", "file.Get").Msg("Fail to get file on mongo.")
		return err
	}
	defer func() {
		downloadStream.Close()
	}()

	fullname := dir + name
	cachePath := path.Dir(fullname)
	if !IsExist(cachePath) {
		err := os.MkdirAll(cachePath, os.ModePerm)
		if err != nil {
			log.Panic().Caller().Err(err).Msgf("Fail to create the CachePath: %s", cachePath)
		}
	}

	fw, err := os.Create(fullname)
	if err != nil {
		log.Error().Caller().Err(err).Msgf("Fail to create file. fullname=%s", fullname)
		return err
	}
	defer fw.Close()
	_, err = io.Copy(fw, downloadStream)
	return err
}

// Del 从数据库删除文件
func Del(name, dir, uid string) error {
	if name == "" {
		return errors.New(keys.ErrorParam)
	}

	bucketOptions := options.GridFSBucket().SetName(mongo.Files)
	bucket, _ := gridfs.NewBucket(mongo.Client.Database(config.App.Mongo.Database), bucketOptions)
	downloadStream, err := bucket.OpenDownloadStreamByName(name)
	if err != nil {
		return err
	}
	defer func() {
		downloadStream.Close()
	}()

	err = bucket.Delete(downloadStream.GetFile().ID)
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "file.Del").Msgf("Fail to delete from GridFS. Name=%s", name)
		return err
	}

	fullname := dir + name
	err = os.Remove(fullname)
	if err != nil {
		log.Warn().Caller().Err(err).Msgf("Fail to delete file. fullname=%s", fullname)
	}
	return nil
}

// ImageThumbnail 生成图片缩略图
func ImageThumbnail(src string, w, h int) (string, error) {
	filename := fmt.Sprintf("%s_%d_%d", src, w, h)

	if IsExist(filename) {
		return filename, nil
	}

	file, err := os.Open(src)
	if err != nil {
		return src, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return src, err
	}

	thumb := resize.Thumbnail(uint(w), uint(h), img, resize.Lanczos3)
	out, err := os.Create(filename)
	if err != nil {
		return src, err
	}
	defer out.Close()

	// Write new image to file.
	err = jpeg.Encode(out, thumb, nil)

	if err != nil {
		return src, err
	}

	return filename, nil
}

// PatchAttr 修改属性
func PatchAttr(fullname string, data bson.M, uid string) (err error) {
	coll := mongo.Client.Database(config.App.Mongo.Database).Collection("fs." + mongo.Files)

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"filename": fullname, "metadata.userId": uid}
	update := bson.M{"$set": bson.M{"metadata": data}}
	_, err = coll.UpdateOne(context.TODO(), filter, update, opts)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "file.PatchAttr").Msgf("Fail to write mongo: data=%v", data)
		err = errors.New(keys.ErrorSave)
	}

	return
}

// GetFileList 获取目录下所有文件
func GetFileList(path string) []bson.M {
	fs, _ := ioutil.ReadDir(path)

	list := make([]bson.M, len(fs))
	for index, file := range fs {
		list[index] = bson.M{}
		if file.IsDir() {
			list[index]["name"] = file.Name()
			list[index]["dir"] = true
			list[index]["list"] = GetFileList(path + file.Name() + "/")
		} else {
			list[index]["name"] = file.Name()
			list[index]["url"] = strings.Replace(path, CachePath, "/image", 1) + file.Name()
		}
	}

	return list
}
