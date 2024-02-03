package mongodb

import (
	"context"
	"os"

	usersColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/users"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongodbURI      string
	mongodbDatabase string

	Client    *mongo.Client
	ClientErr error
)

func Initialize() {
	mongodbURI = os.Getenv("MONGODB_URI")
	mongodbDatabase = os.Getenv("MONGODB_DATABASE")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)

	Client, ClientErr = mongo.Connect(context.TODO(), opts)
	if ClientErr != nil { // client Connection에서 에러가 발생하면
		logger.Error(ClientErr)
		panic(ClientErr)
	}
	logger.Info("MongoDB connected successfully!")

	defineCollection()
	ensureIndexes(context.TODO(), Client.Database(mongodbDatabase))
}

func defineCollection() {
	database := Client.Database(mongodbDatabase)

	usersColl.Define(*database)
}
