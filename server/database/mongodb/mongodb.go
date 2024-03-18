package mongodb

import (
	"context"
	"os"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	loginsColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongodbURI      string
	mongodbDatabase string

	Client *mongo.Client
)

func MustInitialize() {
	mongodbURI = os.Getenv("MONGODB_URI")
	mongodbDatabase = os.Getenv("MONGODB_DATABASE")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)

	var err error
	Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil { // client Connection에서 에러가 발생하면
		logger.Error(err)
		panic(err)
	}
	logger.Info("MongoDB connected successfully!")

	defineCollection()
	MustEnsureIndexes(Client.Database(mongodbDatabase))
}

func defineCollection() {
	database := Client.Database(mongodbDatabase)

	usersColl.Define(*database)
	loginsColl.Define(*database)
}
