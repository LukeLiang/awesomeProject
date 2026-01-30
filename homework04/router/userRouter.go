package router

import "github.com/gin-gonic/gin"

type UserRouter struct {
}

func (r *UserRouter) InitUserRouter(rg *gin.RouterGroup) {
	router := rg.Group("/user")
	{
		router.POST("/register", userApi.Register)
		router.POST("/login", userApi.Login)
	}
}
