package tokensColl

import (
	"context"

	mongo_common "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenRepository interface {
	// [tokens] tokens Collection에서 uid를 통해 Token을 조회하는 함수
	/* ErrMongoInvalidObjectIDError, ErrTokenNotRegistered, ErrMongoFindError */
	FindTokenByUid(uid string) (token TokenSchema, err error)
	// [tokens] tokens Collection에서 refreshToken으로 uid를 포함하는 Token을 조회하는 함수
	/* ErrTokenNotRegistered, ErrMongoFindError */
	FindTokenByRefreshToken(refreshToken string) (token TokenSchema, err error)
	// [tokens] tokens Collection에 토큰을 새로 등록하거나 기존 토큰을 refresh하는 함수.
	/* ErrMongoUpdateError */
	UpsertToken(ctx context.Context, token TokenSchema) (err error)
}

type MongoTokenRepository struct {
	collection *mongo.Collection
}

type TokenSchema struct {
	Uid          primitive.ObjectID `bson:"_id,omitempty"`
	RefreshToken string             `bson:"refreshToken,omitempty"`
	ExpiredAt    primitive.DateTime `bson:"expiredAt,omitempty"`
}

const (
	CollectionName = "tokens"
)

var (
	Collection = &mongo.Collection{}
)

var CollectionInfo mongo_common.CollectionInfo = mongo_common.CollectionInfo{
	Name: CollectionName,
	Indexes: []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "refreshToken", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "expiredAt", Value: -1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	},
}

func NewMongoUserRepository() *MongoTokenRepository {
	return &MongoTokenRepository{
		collection: Collection,
	}
}

func Define(database mongo.Database) {
	Collection = database.Collection(CollectionName)
}
