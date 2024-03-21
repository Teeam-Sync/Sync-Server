package utils_errors

import "errors"

var (
	/* [common] */
	// 예상치 못한 에러 -> 발생한다면 꼭 예외처리가 필요함
	ErrUnexpectedError = errors.New("[common] Unexpected Error Occured")
	// 다루어지지 않은 환경변수로 발생한 에러
	ErrUnhandledEnvironmentVariable = errors.New("[common] Environment Variable Unhandled")
	// 옳바르지 않은 환경변수로 발생한 에러
	ErrInvalidEnvironmentVariable = errors.New("[common] Environment Variable Invalid")

	/* [token] */
	// 유효하지 않은 토큰
	ErrInvalidToken = errors.New("[token] Failed To Verifing Token")
	// 만료된 Access Token
	ErrExpiredAccessToken = errors.New("[token] Access Token Expired")
	// 만료된 Refresh Token
	ErrExpiredRefreshToken = errors.New("[token] Refresh Token Expired")
)
