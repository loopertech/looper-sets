package db

import (
	"context"
	"log"
	"os"

	"github.com/edgedb/edgedb-go"
)

func ConnectDB() (*edgedb.Client, error) {
	log.Println("Connecting to EdgeDB...")
	var options edgedb.Options
	if os.Getenv("APP_ENV") == "dev" {
		options = edgedb.Options{
			Database: "edgedb",
			User: "edgedb",
		}
	} else {
		// Handle remote connection here
	}

	db, err := edgedb.CreateClient(context.Background(), options)
	log.Println("EdgeDB client connected...")
	return db, err
}