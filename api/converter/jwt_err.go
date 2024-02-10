package converter

import "errors"

var (
	ErrUnexpectedSigningMethodError = errors.New("[token] Unexpected Signing Method Error Occured")
	ErrUnverifiableTokenError = errors.New("[token] Unverifiable Token Error Occured")
)