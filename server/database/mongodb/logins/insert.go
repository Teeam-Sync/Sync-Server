package loginsColl

import (
	"context"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// [logins] Login Collection에 User를 넣는 함수
/* ErrUserAlreadyRegistered, ErrMongoInsertError */
func (r *MongoLoginRepository) InsertLoginUser(ctx context.Context, loginUser LoginSchema) (err error) {
	_, err = r.collection.InsertOne(ctx, loginUser)
	if mongo.IsDuplicateKeyError(err) { // 이미 등록된 이메일이 있을때...
		return utils_errors.ErrUserAlreadyRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return utils_errors.ErrMongoInsertError
	}

	return nil
}
