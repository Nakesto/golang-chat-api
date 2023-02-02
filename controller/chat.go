package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nakesto/chat-api/models"
	"github.com/nakesto/chat-api/token"
)

type ChatRoomInput struct {
	Sender   string `json:"sender" form:"sender" binding:"required"`
	Receiver string `json:"receiver" form:"receiver" binding:"required"`
}

type ChatsInput struct {
	Sender   string `json:"sender" form:"sender" binding:"required"`
	Receiver string `json:"receiver" form:"receiver" binding:"required"`
	Message  string `json:"message" form:"message" binding:"required"`
}

func GetChatRoom(c *gin.Context) {

	var room []models.ChatRoom

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		fmt.Println("User not valid")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userId)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.DB.Preload("Sender").Preload("Receiver").Where("sender_name = ?", u.Username).Find(&room).Error

	fmt.Println(room)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"chatroom": room,
	})
}

func GetActiveChat(c *gin.Context) {

	var messages []models.Chat

	receiverName := c.Query("name")

	// if receiverName == "" {

	// }

	if receiverName == "" {
		fmt.Println("request not valid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	fmt.Println(receiverName)

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		fmt.Println("User not valid")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userId)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	set1 := []string{u.Username, receiverName}
	set2 := []string{receiverName, u.Username}

	err = models.DB.Where("(sender_name, receive_name) IN ((?),(?))", set1, set2).Find(&messages).Error

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "OK",
		"chats":   messages,
	})
}

func AddActiveChat(c *gin.Context) {
	var input ChatsInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chats := models.Chat{}

	chats.SenderName = input.Sender
	chats.ReceiveName = input.Receiver
	chats.Message = input.Message

	chat, err := chats.SaveChat()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	set1 := []string{input.Sender, input.Receiver}
	set2 := []string{input.Receiver, input.Sender}

	err = models.DB.Model(&models.ChatRoom{}).Where("(sender_name, receive_name) IN ((?),(?))", set1, set2).Update(&models.ChatRoom{LastMessage: chat.Message}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"chatroom": chat,
	})
}

func AddChatRoom(c *gin.Context) {
	var input ChatRoomInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := models.ChatRoom{}

	room.SenderName = input.Sender
	room.ReceiveName = input.Receiver

	chatroom, err := room.SaveRoom()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room2 := models.ChatRoom{}

	room2.SenderName = input.Receiver
	room2.ReceiveName = input.Sender

	_, err = room2.SaveRoom()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"chatroom": chatroom,
	})
}

func GetFriends(c *gin.Context) {
	params := c.Query("name")

	if params == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Params name required"})
		return
	}

	u, err := models.GetUserByUsername(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"users": u,
	})
}
