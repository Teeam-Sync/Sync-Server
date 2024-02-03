package mongodb

import (
	"context"

	mongo_common "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/common"
	loginsColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/users"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionInfos = []mongo_common.CollectionInfo{
	usersColl.CollectionInfo,
	loginsColl.CollectionInfo,
}

func ensureIndexes(ctx context.Context, db *mongo.Database) {
	for _, info := range collectionInfos {
		_, err := db.Collection(info.Name).Indexes().CreateMany(ctx, info.Indexes)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
	}
}
