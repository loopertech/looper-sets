package main

import (
	"fmt"
	"log"
	edb "looper-sets-backend/pkg/db"
	"looper-sets-backend/pkg/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting backend...")

	// Load env vars
	fmt.Println("Loading environment variables...")
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file...")
	}
	fmt.Println(os.Getenv("APP_ENV"))

	
	db, dbError := edb.ConnectDB()
	if dbError != nil {
		log.Fatal(dbError)
	}


	// Setup server
	host := "localhost:"
	if os.Getenv("APP_ENV") != "dev" {
		gin.SetMode(gin.ReleaseMode)
		host = ":"
	}
	server := gin.Default()

	// Add routes
	routes.PingPong(server)
	routes.Users(server, db)
	routes.Auth(server, db)
 
	// Start server
	log.Println("Starting server...")
	server.Run(host + "8080")
}