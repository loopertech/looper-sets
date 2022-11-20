package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	fmt.Println("Loading environment variables...")

	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file...")
	}
	fmt.Println(os.Getenv("APP_ENV"))
}