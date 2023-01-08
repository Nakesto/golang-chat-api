package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakesto/chat-api/controller"
)

func Protected(protectedRoutes *gin.RouterGroup) {
	protectedRoutes.GET("/current", controller.CurrentUser)
	protectedRoutes.POST("/addchatroom", controller.AddChatRoom)
	protectedRoutes.GET("/chatroom", controller.GetChatRoom)
	protectedRoutes.GET("/activechat", controller.GetActiveChat)
}
