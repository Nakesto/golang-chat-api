package main

import (
	"net/http"
	"os"

	"github.com/nakesto/chat-api/middleware"
	"github.com/nakesto/chat-api/models"
	"github.com/nakesto/chat-api/routes"
	"github.com/nakesto/chat-api/token"

	"github.com/gin-gonic/gin"
)

func main() {
	models.SetupModels()

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	hub := newHub()
	go hub.run()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Content-Type", "application/json")
		c.Next()
	})

	publicRoutes := r.Group("/")
	{
		routes.Public(publicRoutes)
	}

	protectedRoutes := r.Group("/api")
	{
		protectedRoutes.Use(middleware.JwtAuthMiddleware())

		protectedRoutes.GET("/ws", func(c *gin.Context) {
			userID, err := token.ExtractTokenID(c)

			if err != nil {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
			}

			u, err := models.GetUserByID(userID)

			if err != nil {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
			}

			serveWs(hub, c.Writer, c.Request, u.Username)
		})

		routes.Protected(protectedRoutes)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r.Run("127.0.0.1:" + port)
}
