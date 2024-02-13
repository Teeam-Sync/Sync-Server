package utils_jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

func CreateToken(uid string) (authToken string, err error) {
	claim := AuthTokenClaims{
		Uid: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	secretCode := os.Getenv("JWT_SECRET_KEY")
	signedAuthToken, err := token.SignedString([]byte(secretCode))
	if err != nil {
		return "", err
	}
	return signedAuthToken, nil
}
