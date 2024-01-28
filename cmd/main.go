package main

import (
	"log"

	grpc_handler "github.com/Teeam-Sync/Sync-Server/api/handler"
	"github.com/Teeam-Sync/Sync-Server/internal/database"
	"github.com/joho/godotenv"
)

const (
	ENV_FILE = ".env"
)

func main() {
	initialize()
}

func init() {
	defer func() { // panic 발생 시 recover
		if r := recover(); r != nil {
			log.Println("Recovered from panic during initialization:", r)
		}
	}()

	err := godotenv.Load(ENV_FILE)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func initialize() {
	database.Initialize()
	grpc_handler.Initialize()
}
