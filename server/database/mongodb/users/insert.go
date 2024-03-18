package usersColl

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

// [users] users Collection에 새로운 유저를 등록하는 함수.
/* ErrUserAlreadyRegistered, ErrMongoInsertError */
func (r *MongoUserRepository) InsertUser(ctx context.Context, userColl UserSchema) (err error) {
	_, err = r.collection.InsertOne(ctx, userColl)
	if mongo.IsDuplicateKeyError(err) { // 이미 등록된 Email 계정의 User일 경우
		return converter.ErrUserAlreadyRegistered
	} else if err != nil { // 그 이외의 Insert Error
		logger.Error(err)
		return converter.ErrMongoInsertError
	}

	return nil
}
