package tokensColl

import (
	"context"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"github.com/Teeam-Sync/Sync-Server/utils"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// [tokens] tokens Collection에 토큰을 새로 등록하거나 기존 토큰을 refresh하는 함수.
/* ErrMongoUpdateError */
// TODO : 현 방식으로는 여러 기기에서 같이 로그인을 사용하게 될 경우 refreshToken이 가장 최근 기기에서만 작동한다는 단점이 있다.
// 추후에 해당 부분을 device id를 identifier로 추가하는 등의 방법으로 기기마다 refreshToken이 생성될 수 있도록 하는 것이 좋을 것 같다.
func (r *MongoTokenRepository) UpsertToken(ctx context.Context, token TokenSchema) (err error) {
	token.RefreshToken = utils.MakeHash(token.RefreshToken)

	filter := TokenSchema{Uid: token.Uid}
	update := bson.M{"$set": token}
	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil { // unexpected error
		logger.Error(err)
		return utils_errors.ErrMongoUpdateError
	}

	return nil
}
