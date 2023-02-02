package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakesto/chat-api/controller"
)

func Public(publicRoutes *gin.RouterGroup) {
	publicRoutes.GET("/user/:id", controller.FindUserByID)
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)
	publicRoutes.GET("/test", controller.Test)
	publicRoutes.GET("/", func(c *gin.Context){
	c.JSON(200, "Welcome to golang chat api")
	})
}
