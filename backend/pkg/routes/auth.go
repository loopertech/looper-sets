package routes

import (
	"looper-sets-backend/pkg/middleware"
	service "looper-sets-backend/pkg/services"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

func Auth(server *gin.Engine, db *edgedb.Client) {
	authNoToken := server.Group("/auth")
	authWithToken := server.Group("/auth")
	authWithToken.Use(middleware.VerifyUser)
	{
		// Login
		authNoToken.POST("/login", func(c *gin.Context) {
			service.Login(c, db)
		})
		// Login
		authWithToken.POST("/change-password", func(c *gin.Context) {
			service.ChangePassword(c, db)
		})
		// Me
		authWithToken.POST("/me", func(c *gin.Context) {
			service.GetMe(c, db)
		})
	}
}