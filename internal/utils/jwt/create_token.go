package utils_jwt

import (
	"os"
	"time"

	loginsColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/logins"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthTokenClaims struct {
	TokenUUID string `bson:"tid"`
	Uid string `bson:"_id"`
	Email string `bson:"email"`
	jwt.StandardClaims
}

func CreateToken(loginUser loginsColl.LoginsSchema, find_err error) (authToken string, err error) {
	at := AuthTokenClaims{
		TokenUUID: uuid.NewString(),
		Uid: loginUser.Uid.String(),
		Email: loginUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*2).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	secretCode := os.Getenv("JWT_SECRET_KEY")
	signedAuthToken, err := token.SignedString([]byte(secretCode))
	if err != nil {
		return "", err 
	}
	return signedAuthToken, nil 
}

