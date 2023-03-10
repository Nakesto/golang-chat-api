package controller

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cloudinary "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/nakesto/chat-api/models"
	"github.com/nakesto/chat-api/token"
)

type register struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func FindUserByID(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 64)

	var u models.User

	err := models.DB.Model(models.User{}).Where("uid = ?", userId).Take(&u).Error

	if err != nil {
		log.Fatal(err.Error())
	}

	c.JSON(200, gin.H{
		"user": u,
	})
}

func Register(c *gin.Context) {
	var input register

	file, err := c.FormFile("photo")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Photo required for this request"})
		return
	}

	image, err := file.Open()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cloud, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	res, err := cloud.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: file.Filename+input.Username+strconv.FormatInt(time.Now().Unix(), 10),
		Tags: strings.Split("-", file.Filename + input.Username),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password
	u.PhotoURL = res.SecureURL
	u.Photo = res.PublicID

	user, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

type login struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input login

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, token, err := models.LoginCheck(input.Username, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "Your account has logined",
		"user":   u,
		"token":  token,
	})
}

func CurrentUser(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func Test(c *gin.Context) {

	db := models.DB

	chats := models.Chat{}

	err := db.Model(&models.User{}).Preload("Sender").Preload("Receiver").First(&chats).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chats})
}

func ChangeProfile(c *gin.Context) {

	file, err := c.FormFile("photo")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Photo required for this request"})
		return
	}

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := file.Open()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cloud, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	res, err := cloud.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: file.Filename+u.Username+strconv.FormatInt(time.Now().Unix(), 10),
		Tags: strings.Split("-", file.Filename + u.Username),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.UpdatePhoto(models.User{
		Username: u.Username,
		Photo: res.PublicID,
		PhotoURL: res.SecureURL,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := cloud.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: u.Photo,
	})

	if err != nil {
  		c.JSON(http.StatusBadRequest, gin.H{"error": "File could not be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": result.Result, "photoURL" : res.SecureURL})
}

type status struct{
	Status string `json:"status" form:"status" binding:"required"`
}

func ChangeStatus(c *gin.Context){
	var input status

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}

	user.Username = u.Username
	user.Status = input.Status

	err = models.UpdateStatus(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Status has been changed"})
}
