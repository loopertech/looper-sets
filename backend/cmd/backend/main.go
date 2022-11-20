package main

import (
	"log"
	"looper-sets-backend/pkg/config"
	"looper-sets-backend/pkg/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting backend...")

	config.LoadEnvironmentVariables()

	host := "localhost:"
	if os.Getenv("APP_ENV") != "dev" {
		gin.SetMode(gin.ReleaseMode)
		host = ":"
	}

	server := gin.Default()
	routes.AddPingPongRoutes(server)
 
	log.Println("Starting server...")
	server.Run(host + "8080")
}