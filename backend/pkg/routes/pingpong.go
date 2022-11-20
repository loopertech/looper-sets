package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPingPongRoutes(server *gin.Engine) {
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})
}