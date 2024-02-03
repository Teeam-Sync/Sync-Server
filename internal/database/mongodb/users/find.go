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
	if err == mongo.ErrNoDocuments { // 해당하는 email에 User가 등록되지 않은 경우
		return userColl, converter.ErrUserNotRegistered
	} else if err != nil { // 그 이외의 Find Error
		logger.Error(err)
		return userColl, converter.ErrMongoFindError
	}

	return userColl, nil
}

// [users] Uid를 통해서 User를 가져오는 함수 // ErrUserNotRegistered, ErrMongoFindError
func FindLoginUserByUid(uid string) (user UsersSchema, err error) {
	filter := bson.M{"_id": uid}

	err = Collection.FindOne(context.TODO(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments { // Collection에서 User가 등록되어있지 않는 경우
		return user, converter.ErrUserNotRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return user, converter.ErrMongoFindError
	}

	return user, nil
}
