package loginsColl

import (
	"context"

	mongo_common "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginRepository interface {
	// [logins] Email을 통해서 Login 유저를 가져오는 함수
	/* ErrUserNotRegistered, ErrMongoFindError */
	FindLoginUserByEmail(email string) (loginUser LoginSchema, err error)
	// [logins] Uid를 통해서 Login 유저를 가져오는 함수
	/* ErrUserNotRegistered, ErrMongoFindError */
	FindLoginUserByUid(uid string) (loginUser LoginSchema, err error)
	// [logins] Login Collection에 User를 넣는 함수
	/* ErrUserAlreadyRegistered, ErrMongoInsertError */
	InsertLoginUser(ctx context.Context, loginUser LoginSchema) (err error)
}

type MongoLoginRepository struct {
	collection *mongo.Collection
}

type LoginSchema struct {
	Uid       primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email,omitempty"`
	Password  string             `bson:"password,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt,omitempty"`
}

const (
	CollectionName = "logins"
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
	},
}

func NewMongoLoginRepository() *MongoLoginRepository {
	return &MongoLoginRepository{
		collection: Collection,
	}
}

func Define(database mongo.Database) {
	Collection = database.Collection(CollectionName)
}
