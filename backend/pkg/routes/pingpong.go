package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingPong(server *gin.Engine) {
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})
}