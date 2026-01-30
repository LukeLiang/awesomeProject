package router

import (
	"fmt"

	"awesomeProject/homework04/global"
	"awesomeProject/homework04/middleware"

	"github.com/gin-gonic/gin"
)

// InitServer 初始化路由
func InitServer() {
	router := gin.New()
	InitRouter(router)

	// 启动服务
	_ = router.Run(fmt.Sprintf(":%d", global.Config.Server.Port))

}

// InitRouter 初始化路由
func InitRouter(router *gin.Engine) {
	routers := RoutersGroup{}

	// 注册全局中间件
	router.Use(middleware.CorsHandler())   // 跨域处理
	router.Use(middleware.ErrorHandler())  // panic 恢复
	router.Use(middleware.RequestLogger()) // 请求日志

	// 注册 404 和 405 处理器
	router.NoRoute(middleware.NotFoundHandler())
	router.NoMethod(middleware.MethodNotAllowedHandler())

	// 非鉴权 api
	publicRouter := router.Group(fmt.Sprintf("/%s/%s", global.Config.Prefix, global.Config.Version))
	{
		routers.InitUserRouter(publicRouter)
	}

	// 鉴权 api
	privateRouter := router.Group(fmt.Sprintf("/%s/%s", global.Config.Prefix, global.Config.Version))
	privateRouter.Use(middleware.VerifyToken())
	{
		routers.initPostRouter(privateRouter)
	}
}
