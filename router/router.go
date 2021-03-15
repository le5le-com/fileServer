package router

import (
	"strconv"

	"fileServer/config"
	"fileServer/file"
	"fileServer/middlewares"

	"github.com/kataras/iris/v12"
)

// Listen 监听路由
func Listen() {
	route := iris.New()
	route.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.JSON("Internal server in file server!")
	})

	route.Use(middlewares.Usr)

	// 文件模块
	file.Route(route)

	// 监听
	port := strconv.Itoa(int(config.App.Port))
	route.Run(
		iris.Addr(":"+port),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}
