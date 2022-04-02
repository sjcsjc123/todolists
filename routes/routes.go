package routes

import (
	"TodoLists/api"
	"TodoLists/middleware"
	"github.com/gin-gonic/gin"
)

func NewRoutes() *gin.Engine {
	r := gin.Default()
	r.POST("/login", api.Login)
	r.POST("/register", api.Register)
	group := r.Group("/", middleware.ValidToken())
	{
		group.GET("/list", api.ListTask)
		group.POST("/create", api.CreateTask)
		group.GET("/finish", api.FinishTask)
		group.GET("/delete", api.DeleteTask)
		routerGroup := group.Group("/list")
		{
			//1代表已完成，0代表未完成
			routerGroup.GET("/noFinish", api.ListNoFinishTask)
			routerGroup.GET("/finish", api.ListFinishTask)
		}
	}
	return r
}
