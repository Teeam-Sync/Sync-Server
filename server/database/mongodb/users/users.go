package usersColl

import (
	"context"

	mongo_common "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	FindUserByEmail(email string) (userColl UserSchema, err error)
	FindLoginUserByUid(uid string) (user UserSchema, err error)
	InsertUser(ctx context.Context, userColl UserSchema) (err error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

type UserSchema struct {
	Uid       primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Email     string             `bson:"email,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt,omitempty"`
}

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

func NewMongoUserRepository() *MongoUserRepository {
	return &MongoUserRepository{
		collection: Collection,
	}
}

func Define(database mongo.Database) {
	Collection = database.Collection(CollectionName)
}
