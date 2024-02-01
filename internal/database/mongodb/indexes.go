package mongodb

import (
	"context"

	mongo_common "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/common"
	usersColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/users"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionInfos = []mongo_common.CollectionInfo{
	usersColl.CollectionInfo,
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	for _, info := range collectionInfos {
		_, err := db.Collection(info.Name).Indexes().CreateMany(ctx, info.Indexes)
		if err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}
