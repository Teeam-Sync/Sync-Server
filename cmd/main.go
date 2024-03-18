package main

import (
	"log"

	database "github.com/Teeam-Sync/Sync-Server/server/database/mongodb"
	grpc_handler "github.com/Teeam-Sync/Sync-Server/server/grpc"
	utils_kst "github.com/Teeam-Sync/Sync-Server/utils/kst"
	"github.com/joho/godotenv"
)

const (
	ENV_FILE = ".env"
)

func main() {
	mustInitialize()
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

func mustInitialize() {
	utils_kst.MustLoadKST()
	database.MustInitialize()
	grpc_handler.Initialize()
}
