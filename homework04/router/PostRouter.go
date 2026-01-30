package router

import (
	"github.com/gin-gonic/gin"
)

type PostRouter struct {
}

func (r *PostRouter) initPostRouter(rg *gin.RouterGroup) {
	router := rg.Group("/post")

	{
		router.POST("/create", postApi.Create)
		router.GET("/list/:userId", postApi.ListPosts)
		router.GET("/detail/:id", postApi.Detail)
		router.PUT("/update/:id", postApi.Update)
		router.DELETE("/delete/:id", postApi.Delete)
	}
}
