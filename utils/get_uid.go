package utils

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
)

// context(Metadata & JWT Token)로부터 uid를 가져오는 함수
/* ErrUnexpectedError */
func GetUid(ctx context.Context) (uid string, err error) {
	uid, ok := ctx.Value(converter.UidKey).(string)
	if !ok { // unexpected error
		logger.Error(err)
		return "", utils_errors.ErrUnexpectedError
	}

	return uid, err
}

// context(Metadata)로부터 JWT refresh token을 가져오는 함수
/* ErrUnexpectedError */
func GetRefreshToken(ctx context.Context) (refreshToken string, err error) {
	uid, ok := ctx.Value(converter.RefreshTokenKey).(string)
	if !ok { // unexpected error
		logger.Error(err)
		return "", utils_errors.ErrUnexpectedError
	}

	return uid, err
}
