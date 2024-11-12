package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register routes to the Gin router
func RegisterRoutes(router *gin.Engine) {
	router.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
