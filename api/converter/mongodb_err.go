package converter

import "errors"

var (
	// [common]
	ErrUnexpectedError = errors.New("[common] Unexpected Error Occured")

	// [mongo]
	ErrMongoInsertError = errors.New("[mongo] MongoDB Insert Function Occured Error")
	ErrMongoFindError   = errors.New("[mongo] MongoDB Find Function Occured Error")
	ErrMongoUpdateError = errors.New("[mongo] MongoDB Update Function Occured Error")
	ErrMongoDeleteError = errors.New("[mongo] MongoDB Delete Function Occured Error")

	// [users]
	ErrUserNotRegistered     = errors.New("[users] User Not Registered Before")
	ErrUserAlreadyRegistered = errors.New("[users] User Already Registered Before")
)
