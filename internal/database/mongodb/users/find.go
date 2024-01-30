package usersColl

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// [users] users Collection에서 Email을 통해 User를 조회하는 함수 // ErrMongoFindError, ErrMongoFindError
func FindUserByEmail(email string) (userColl UsersSchema, err error) {
	filter := bson.M{
		"email": email,
	}

	err = Collection.FindOne(context.TODO(), filter).Decode(&userColl)
	if err != mongo.ErrNoDocuments { // 해당하는 email에 User가 등록되지 않은 경우
		return userColl, converter.ErrUserNotRegistered
	} else if err != nil { // 그 이외의 Find Error
		logger.Error(err)
		return userColl, converter.ErrMongoFindError
	}

	return userColl, nil
}
