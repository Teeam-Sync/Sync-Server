package mongo_common

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionInfo struct {
	Name    string
	Indexes []mongo.IndexModel
}
