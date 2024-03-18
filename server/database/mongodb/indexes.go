package mongodb

import (
	"context"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	mongo_common "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/common"
	loginsColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/users"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionInfos = []mongo_common.CollectionInfo{
	usersColl.CollectionInfo,
	loginsColl.CollectionInfo,
}

func MustEnsureIndexes(db *mongo.Database) {
	for _, info := range collectionInfos {
		_, err := db.Collection(info.Name).Indexes().CreateMany(context.Background(), info.Indexes)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
	}
}
