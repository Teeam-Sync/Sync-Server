package usersColl

import (
	"context"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// [users] users Collection에서 Email을 통해 User를 조회하는 함수
/* ErrUserNotRegistered, ErrMongoFindError */
func (r *MongoUserRepository) FindUserByEmail(email string) (user UserSchema, err error) {
	filter := UserSchema{Email: email}

	err = r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments { // 해당하는 email에 User가 등록되지 않은 경우
		return user, utils_errors.ErrUserNotRegistered
	} else if err != nil { // 그 이외의 Find Error
		logger.Error(err)
		return user, utils_errors.ErrMongoFindError
	}

	return user, nil
}

// [users] Uid를 통해서 User를 가져오는 함수
/* ErrMongoInvalidObjectIDError, ErrUserNotRegistered, ErrMongoFindError */
func (r *MongoUserRepository) FindLoginUserByUid(uid string) (user UserSchema, err error) {
	parsedUid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return user, utils_errors.ErrMongoInvalidObjectIDError
	}

	filter := UserSchema{Uid: parsedUid}

	err = r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments { // Collection에서 User가 등록되어있지 않는 경우
		return user, utils_errors.ErrUserNotRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return user, utils_errors.ErrMongoFindError
	}

	return user, nil
}
