package loginsColl

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// [logins] Email을 통해서 Login 유저를 가져오는 함수 // ErrUserNotRegistered, ErrMongoFindError
func FindLoginUserByEmail(email string) (loginUser LoginsSchema, err error) {
	filter := LoginsSchema{Email: email}

	err = Collection.FindOne(context.TODO(), filter).Decode(&loginUser)
	if err == mongo.ErrNoDocuments { // login Collection에서 User가 등록되어있지 않는 경우
		return loginUser, converter.ErrUserNotRegistered
	} else if err != nil { // unexpected error
		return loginUser, converter.ErrMongoFindError
	}

	return loginUser, nil
}

// [logins] Uid를 통해서 Login 유저를 가져오는 함수 // ErrUserNotRegistered, ErrMongoFindError
func FindLoginUserByUid(uid string) (loginUser LoginsSchema, err error) {
	parsedUid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return loginUser, converter.ErrMongoInvalidObjectIDError
	}

	filter := LoginsSchema{Uid: parsedUid}

	err = Collection.FindOne(context.TODO(), filter).Decode(&loginUser)
	if err == mongo.ErrNoDocuments { // login Collection에서 User가 등록되어있지 않는 경우
		return loginUser, converter.ErrUserNotRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return loginUser, converter.ErrMongoFindError
	}

	return loginUser, nil
}
