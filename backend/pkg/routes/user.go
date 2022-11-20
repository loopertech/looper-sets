package routes

import (
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

func Users(server *gin.Engine, db *edgedb.Client) {
	users := server.Group("/users")
	// users.Use()
	{
		// Create user
		users.POST("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "create user")
		})
		// Get all users
		users.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "get users")
		})
		// Get user
		users.GET("/:uuid", func(c *gin.Context) {
			c.JSON(http.StatusOK, "get user")
		})
		// Edit user
		users.PATCH("/:uuid", func(c *gin.Context) {
			c.JSON(http.StatusOK, "edit user")
		})
		// Delete user
		users.DELETE("/:uuid", func(c *gin.Context) {
			c.JSON(http.StatusOK, "delete user")
		})
	}
}