package loginsColl

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	"go.mongodb.org/mongo-driver/mongo"
)

// [logins] Login Collection에 User를 넣는 함수 // ErrUserAlreadyRegistered, ErrMongoInsertError
func InsertLoginUser(ctx context.Context, loginUser LoginsSchema) (err error) {
	_, err = Collection.InsertOne(ctx, loginUser)
	if mongo.IsDuplicateKeyError(err) { // 이미 등록된 이메일이 있을때...
		return converter.ErrUserAlreadyRegistered
	} else if err != nil { // unexpected error
		return converter.ErrMongoInsertError
	}

	return nil
}
