package utils_jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// type CustomClaims struct {
// 	Uid string `json:"uid"`
// 	jwt.RegisteredClaims
// }

func CreateToken(uid string) (authToken string, err error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer: "sync",
		Subject: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAuthToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedAuthToken, nil
}