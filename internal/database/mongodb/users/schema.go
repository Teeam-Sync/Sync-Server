package usersColl

import (
	mongo_common "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionName = "users"
)

var (
	Collection = &mongo.Collection{}
)

var CollectionInfo mongo_common.CollectionInfo = mongo_common.CollectionInfo{
	Name: CollectionName,
	Indexes: []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	},
}

type UsersSchema struct {
	Uid       primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Email     string             `bson:"email,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt,omitempty"`
}

func Define(database mongo.Database) {
	Collection = database.Collection(CollectionName)
}
