package usersColl

import (
	"context"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// [users] users Collection에 새로운 유저를 등록하는 함수.
/* ErrUserAlreadyRegistered, ErrMongoInsertError */
func (r *MongoUserRepository) InsertUser(ctx context.Context, user UserSchema) (err error) {
	_, err = r.collection.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) { // 이미 등록된 Email 계정의 User일 경우
		return utils_errors.ErrUserAlreadyRegistered
	} else if err != nil { // 그 이외의 Insert Error
		logger.Error(err)
		return utils_errors.ErrMongoInsertError
	}

	return nil
}
