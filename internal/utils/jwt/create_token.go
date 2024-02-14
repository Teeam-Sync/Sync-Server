package utils_jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func CreateToken(uid string) (authToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Uid": uid,
		"exp": time.Now().Add(time.Hour*2).Unix(),
	})
	signedAuthToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedAuthToken, nil
}