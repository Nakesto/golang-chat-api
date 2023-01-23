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

	if os.Getenv("type") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	hub := newHub()
	go hub.run()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "DELETE, POST, GET, PUT, OPTIONS")
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		} else {
			c.Next()
		}
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

	r.Run(":" + port)
}
