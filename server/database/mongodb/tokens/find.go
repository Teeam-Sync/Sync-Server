package tokensColl

import (
	"context"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"github.com/Teeam-Sync/Sync-Server/utils"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// [tokens] tokens Collection에서 uid를 통해 Token을 조회하는 함수
/* ErrMongoInvalidObjectIDError, ErrTokenNotRegistered, ErrMongoFindError */
func (r *MongoTokenRepository) FindTokenByUid(uid string) (token TokenSchema, err error) {
	parsedUid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return token, utils_errors.ErrMongoInvalidObjectIDError
	}

	filter := TokenSchema{Uid: parsedUid}

	err = r.collection.FindOne(context.Background(), filter).Decode(&token)
	if err == mongo.ErrNoDocuments { // 해당하는 uid에 Token이 없는 경우
		return token, utils_errors.ErrTokenNotRegistered
	} else if err != nil { // 그 이외의 Find Error
		logger.Error(err)
		return token, utils_errors.ErrMongoFindError
	}

	return token, nil
}

// [tokens] tokens Collection에서 refreshToken으로 uid를 포함하는 Token을 조회하는 함수
/* ErrTokenNotRegistered, ErrMongoFindError */
func (r *MongoTokenRepository) FindTokenByRefreshToken(refreshToken string) (token TokenSchema, err error) {
	hashedRefreshToken := utils.MakeHash(refreshToken)

	filter := TokenSchema{RefreshToken: hashedRefreshToken}

	err = r.collection.FindOne(context.Background(), filter).Decode(&token)
	if err == mongo.ErrNoDocuments { // 부합하는 refreshToken을 갖고 있는 Token이 없는 경우
		return token, utils_errors.ErrTokenNotRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return token, utils_errors.ErrMongoFindError
	}

	return token, nil
}
