package routes

import (
	"looper-sets-backend/pkg/middleware"
	service "looper-sets-backend/pkg/services"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

func Users(server *gin.Engine, db *edgedb.Client) {
	users := server.Group("/users")
	usersWithToken := server.Group("/users")
	usersWithToken.Use(middleware.VerifyUser)
	{
		// Create user
		users.POST("/", func(c *gin.Context) {
			service.CreateUser(c, db)
		})
		// Get all users
		users.GET("/", func(c *gin.Context) {
			service.GetUsers(c, db)
		})
		// Get user
		users.GET("/:uuid", func(c *gin.Context) {
			service.GetUser(c, "", db)
		})
		// Edit user
		usersWithToken.PATCH("/:uuid", func(c *gin.Context) {
			service.EditUser(c, db)
		})
		// Delete user
		usersWithToken.DELETE("/:uuid", func(c *gin.Context) {
			service.DeleteUser(c, db)
		})
	}
}