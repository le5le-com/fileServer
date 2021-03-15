package middlewares

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"

	"fileServer/config"
	"fileServer/keys"
	"fileServer/utils"
)

// Usr 解析用户身份
func Usr(ctx iris.Context) {
	if config.App.Jwt == "" {
		ctx.Values().Set("uid", "system")
		ctx.Values().Set("username", "乐乐")
		ctx.Values().Set("role", "operation")
		ctx.Next()
		return
	}

	// 获取header
	data := ctx.GetHeader("Authorization")
	if data == "" {
		data = ctx.URLParam("Authorization")
	}
	if data == "" {
		ctx.Next()
		return
	}

	// jwt校验
	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名方法错误: %v", token.Header["alg"])
		}
		return []byte(config.App.Jwt), nil
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("func", "middlewares.Usr").
			Str("token", data).
			Str("jwt", config.App.Jwt).
			Str("remoteAddr", ctx.RemoteAddr()).
			Msg("Jwt parse error.")
		ctx.Next()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Warn().
			Str("func", "middlewares.Usr").
			Str("token", data).
			Str("jwt", config.App.Jwt).
			Str("remoteAddr", ctx.RemoteAddr()).
			Msg("Jwt invalid.")
		ctx.Next()
		return
	}

	// 设置uid和role
	uid := utils.String(claims["uid"])
	if uid != "" {
		ctx.Values().Set("uid", uid)
		ctx.Values().Set("username", utils.String(claims["username"]))
		ctx.Values().Set("role", utils.String(claims["role"]))
	}

	ctx.Next()
}

// Auth 身份认证中间件
func Auth(ctx iris.Context) {
	if ctx.Values().GetString("uid") != "" {
		ctx.Next()
		return
	}

	ctx.StatusCode(iris.StatusUnauthorized)
	ret := make(map[string]interface{})
	ret["error"] = keys.ErrorNeedSign
	ctx.JSON(ret)
}
