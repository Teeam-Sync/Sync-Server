package utils_jwt

import (
	"os"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"github.com/golang-jwt/jwt"
)

func CheckToken(jwt_token string) (err error) {
	claims := AuthTokenClaims{}
	key := func(jwt_token *jwt.Token) (interface{}, error) {
		if _, ok := jwt_token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, converter.ErrUnexpectedSigningMethodError
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	}

	token, err := jwt.ParseWithClaims(jwt_token, &claims, key)

	if err != nil {
		return converter.ErrUnverifiableTokenError
	} 

	uuid := claims.Uid
	logger.Debug(token.Valid)
	logger.Debug(uuid)
	return nil 
}