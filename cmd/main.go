package main

import (
	"log"
	"os"

	database "github.com/Teeam-Sync/Sync-Server/server/database/mongodb"
	grpc_handler "github.com/Teeam-Sync/Sync-Server/server/grpc"
	jwtService "github.com/Teeam-Sync/Sync-Server/server/service/jwt"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	utils_kst "github.com/Teeam-Sync/Sync-Server/utils/kst"
	"github.com/joho/godotenv"
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

	var envFile string
	appEnv := os.Getenv("APP_ENV")
	switch appEnv {
	case "dev":
		envFile = ".dev.env"
	case "prod":
		envFile = ".prod.env"
	default:
		panic(utils_errors.ErrInvalidEnvironmentVariable)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		panic(utils_errors.ErrInvalidEnvironmentVariable)
	}
}

func mustInitialize() {
	// logger.MustInitialize()
	utils_kst.MustInitialize()
	jwtService.MustInitialize()
	database.MustInitialize()
	grpc_handler.Initialize()
}
