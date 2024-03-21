package jwtService

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	tokensColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/tokens"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type issuerType string

const (
	accessTokenIssuer  issuerType = "sync_access"
	refreshTokenIssuer issuerType = "sync_refresh"
)

var (
	accessTokenKey      []byte
	refreshTokenKey     []byte
	expiredDurationHour time.Duration // for accessToken
	expiredDurationDay  time.Duration // for refreshToken
)

func MustInitialize() {
	jwtAccessTokenKey := os.Getenv("JWT_ACCESSTOKEN_KEY")
	jwtRefreshTokenKey := os.Getenv("JWT_REFRESHTOKEN_KEY")
	jwtExpirationHour := os.Getenv("JWT_ACCESSTOKEN_EXPIRATION")
	jwtExpirationDay := os.Getenv("JWT_REFRESHTOKEN_EXPIRATION")

	if jwtAccessTokenKey == "" || jwtRefreshTokenKey == "" || jwtExpirationHour == "" || jwtExpirationDay == "" {
		logger.Error(utils_errors.ErrUnhandledEnvironmentVariable)
		panic(utils_errors.ErrUnhandledEnvironmentVariable)
	}

	parsedJwtExpirationHour, err := strconv.Atoi(jwtExpirationHour)
	if err != nil {
		logger.Error(err)
		panic(utils_errors.ErrInvalidEnvironmentVariable)
	}

	parsedJwtExpirationDay, err := strconv.Atoi(jwtExpirationDay)
	if err != nil {
		logger.Error(err)
		panic(utils_errors.ErrInvalidEnvironmentVariable)
	}

	accessTokenKey = []byte(jwtAccessTokenKey)
	refreshTokenKey = []byte(jwtRefreshTokenKey)
	expiredDurationHour = time.Duration(parsedJwtExpirationHour) * time.Hour
	expiredDurationDay = time.Duration(parsedJwtExpirationDay) * time.Hour * 24
}

// Createing Token(AccessToken & RefreshToken)
/* ErrMongoInvalidObjectIDError, ErrUnexpectedError */
func CreateJWTToken(uid string) (jwtToken converter.JWTToken, err error) {
	tokenRepository := tokensColl.NewMongoUserRepository()

	parsedUid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return jwtToken, utils_errors.ErrMongoInvalidObjectIDError
	}

	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Subject:  uid,
	}

	accessTokenClaims := claims
	accessTokenClaims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(expiredDurationHour))
	accessTokenClaims.Issuer = accessTokenIssuer.string()

	refreshTokenClaims := claims
	refreshTokenClaims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(expiredDurationDay))
	refreshTokenClaims.Issuer = refreshTokenIssuer.string()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	jwtToken.AccessToken, err = token.SignedString(accessTokenKey)
	if err != nil { // unexpected error
		logger.Error(err)
		return jwtToken, utils_errors.ErrUnexpectedError
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	jwtToken.RefreshToken, err = token.SignedString(refreshTokenKey)
	if err != nil { // unexpected error
		logger.Error(err)
		return jwtToken, utils_errors.ErrUnexpectedError
	}

	tokenRepository.UpsertToken(context.Background(), tokensColl.TokenSchema{
		Uid:          parsedUid,
		RefreshToken: jwtToken.RefreshToken,
		ExpiredAt:    primitive.NewDateTimeFromTime(time.Now().UTC().Add(expiredDurationDay)),
	})

	return jwtToken, nil
}

// Verify JWT Access Token
/* ErrExpiredAccessToken, ErrInvalidToken */
func VerifyAccessToken(tokenString string) (uid string, err error) {
	if tokenString == "" {
		return "", utils_errors.ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok { // cannot verify token
			return nil, utils_errors.ErrInvalidToken
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil { // unexpected error
			return nil, utils_errors.ErrUnexpectedError
		}

		if issuer == accessTokenIssuer.string() { // access token
			return accessTokenKey, nil
		}

		return nil, utils_errors.ErrInvalidToken
	})

	if err != nil || !token.Valid {
		if expiredAt, err := token.Claims.GetExpirationTime(); err == nil {
			if time.Now().UTC().After(expiredAt.Time) {
				return "", utils_errors.ErrExpiredAccessToken
			}
		}
		return "", utils_errors.ErrInvalidToken
	}

	uid, err = token.Claims.GetSubject()
	if err != nil {
		return "", utils_errors.ErrInvalidToken
	}

	return uid, nil
}

// Veify JWT Refresh Token
/* ErrExpiredRefreshToken, ErrUnexpectedError, ErrTokenNotRegistered, ErrInvalidToken */
func VerifyRefreshToken(tokenString string) (uid string, err error) {
	tokenRepository := tokensColl.NewMongoUserRepository()

	refreshToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok { // cannot verify token
			return nil, utils_errors.ErrInvalidToken
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil { // unexpected error
			return nil, utils_errors.ErrUnexpectedError
		}

		if issuer == refreshTokenIssuer.string() { // refresh token
			return refreshTokenKey, nil
		}

		return nil, utils_errors.ErrInvalidToken
	})

	if err != nil || !refreshToken.Valid {
		if expiredAt, err := refreshToken.Claims.GetExpirationTime(); err == nil {
			if time.Now().UTC().After(expiredAt.Time) {
				return "", utils_errors.ErrExpiredRefreshToken
			}
		}
		return "", utils_errors.ErrInvalidToken
	}

	uid, err = refreshToken.Claims.GetSubject()
	if err != nil {
		return "", utils_errors.ErrInvalidToken
	}

	token, err := tokenRepository.FindTokenByRefreshToken(tokenString)
	if err == utils_errors.ErrTokenNotRegistered { // token not registered
		return "", utils_errors.ErrTokenNotRegistered
	} else if err != nil { // unexpected error
		return "", utils_errors.ErrUnexpectedError
	}

	if token.Uid.Hex() != uid { // invalid token
		return "", utils_errors.ErrInvalidToken
	}

	return uid, nil
}

func (i issuerType) string() string {
	return string(i)
}
