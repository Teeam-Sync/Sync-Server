package utils_jwt

import (
	"github.com/Teeam-Sync/Sync-Server/api/converter"
	"github.com/golang-jwt/jwt/v5"
)

func CheckToken(tokenString string) (uid string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil 
	})

	if err != nil {
		return "", err 
	}

	if !token.Valid {
		return "", converter.ErrInvalidTokenError
	}

	claims := token.Claims.(jwt.MapClaims)
	uid = claims["Uid"].(string)

	return uid, nil
}