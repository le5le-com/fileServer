package file

import (
	"fileServer/middlewares"

	"github.com/kataras/iris/v12"
)

// Route file模块路由
func Route(route *iris.Application) {

	route.Get("/file/{path:path}", File)
	route.Post("/file", middlewares.Auth, Limit, UploadFile)
	route.Patch("/file/{path:path}", middlewares.Auth, FileAttr)
	route.Delete("/file/{path:path}", middlewares.Auth, DelFile)

	route.Post("/api/file", middlewares.Auth, Limit, UploadFile)
	route.Patch("/api/file/{path:path}", middlewares.Auth, FileAttr)
	route.Delete("/api/file/{path:path}", middlewares.Auth, DelFile)

	route.Get("/image/{path:path}", Image)
	route.Post("/image", middlewares.Auth, Limit, UploadImage)
	route.Patch("/image/{path:path}", middlewares.Auth, FileAttr)
	route.Delete("/image/{path:path}", middlewares.Auth, DelFile)

	route.Post("/api/image", middlewares.Auth, Limit, UploadImage)
	route.Patch("/api/image/{path:path}", middlewares.Auth, FileAttr)
	route.Delete("/api/image/{path:path}", middlewares.Auth, DelFile)

	route.Get("/images", Images)
}
