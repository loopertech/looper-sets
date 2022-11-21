package routes

import (
	"looper-sets-backend/pkg/middleware"
	service "looper-sets-backend/pkg/services"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

func Songs(server *gin.Engine, db *edgedb.Client) {
	songs := server.Group("/songs")
	songsWithToken := server.Group("/songs")
	songsWithToken.Use(middleware.VerifyUser)
	{
		// Create song
		songsWithToken.POST("/", func(c *gin.Context) {
			service.CreateSong(c, db)
		})
		// Get all songs
		songs.GET("/", func(c *gin.Context) {
			service.GetSongs(c, db)
		})
		// Get song
		songs.GET("/:uuid", func(c *gin.Context) {
			service.GetSong(c, db)
		})
		// // Edit song
		// songsWithToken.PATCH("/:uuid", func(c *gin.Context) {
		// 	service.EditSong(c, db)
		// })
		// Delete song
		songsWithToken.DELETE("/:uuid", func(c *gin.Context) {
			service.DeleteSong(c, db)
		})
	}
}