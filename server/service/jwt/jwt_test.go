package jwtService_test

import (
	"os"
	"testing"

	jwtService "github.com/Teeam-Sync/Sync-Server/server/service/jwt"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"github.com/stretchr/testify/assert"
)

func TestCheckAccessToken(t *testing.T) {
	os.Setenv("JWT_ACCESSTOKEN_KEY", "testAccessTokenKey")
	os.Setenv("JWT_REFRESHTOKEN_KEY", "testRefreshTokenKey")
	os.Setenv("JWT_ACCESSTOKEN_EXPIRATION", "3")  // 3 hours
	os.Setenv("JWT_REFRESHTOKEN_EXPIRATION", "1") // 1 day

	jwtService.MustInitialize()

	expectedUid := "a1b2c3d4"
	t.Run("verified", func(t *testing.T) {
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzeW5jX2FjY2VzcyIsInN1YiI6ImExYjJjM2Q0IiwiZXhwIjo1MzEwOTg1OTczLCJpYXQiOjE3MTA5ODU5NzN9.du60L1tepqrr0INmX3BoHZKYf3T4NHg8IwTxjI-N88Q"

		uid, err := jwtService.VerifyAccessToken(accessToken)
		assert.Nil(t, err)
		assert.Equal(t, expectedUid, uid)
	})

	t.Run("expired Token", func(t *testing.T) {
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzeW5jX2FjY2VzcyIsInN1YiI6ImExYjJjM2Q0IiwiZXhwIjoxNzEwOTg2NTI4LCJpYXQiOjE3MTA5ODY1Mjh9.4w6Rl48u2dKfgxP5hRp3xPk0RALXdLiOBQJwASJNFFI"

		_, err := jwtService.VerifyAccessToken(accessToken)
		assert.EqualError(t, err, utils_errors.ErrExpiredAccessToken.Error())
	})

	t.Run("invalid Token: wrong signature", func(t *testing.T) {
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzeW5jX2FjY2VzcyIsInN1YiI6ImExYjJjM2Q0IiwiZXhwIjo1MzEwOTg1OTczLCJpYXQiOjE3MTA5ODU5NzN9.si-b9ozYo41ulYoFB9Y6Sucov5mMP5V-GOcwtI7LANo"

		_, err := jwtService.VerifyAccessToken(accessToken)
		assert.EqualError(t, err, utils_errors.ErrInvalidToken.Error())
	})

	t.Run("invalid Token: refresh token came", func(t *testing.T) {
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzeW5jX3JlZnJlc2giLCJzdWIiOiJhMWIyYzNkNCIsImV4cCI6NTMxMDk4NTk3MywiaWF0IjoxNzEwOTg1OTczfQ.1AyUPZZ7h7RI5PqnsjmY0b7XgtG4DcgJ7lzRjVqyDo4"

		_, err := jwtService.VerifyAccessToken(accessToken)
		assert.EqualError(t, err, utils_errors.ErrInvalidToken.Error())
	})

	t.Run("invalid Token: refresh token(issuer modulated) came", func(t *testing.T) {
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzeW5jX2FjY2VzcyIsInN1YiI6ImExYjJjM2Q0IiwiZXhwIjo1MzEwOTg1OTczLCJpYXQiOjE3MTA5ODU5NzN9.0YBfGjK6f1yvJ1drgaxD8BmyP0jLjYNbKPh_3R0xZDA"

		_, err := jwtService.VerifyAccessToken(accessToken)
		assert.EqualError(t, err, utils_errors.ErrInvalidToken.Error())
	})

	t.Run("invalid Token: empty token", func(t *testing.T) {
		accessToken := ""

		_, err := jwtService.VerifyAccessToken(accessToken)
		assert.EqualError(t, err, utils_errors.ErrInvalidToken.Error())
	})
}
