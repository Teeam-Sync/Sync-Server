package utils_errors

import "errors"

var (
	// [mongo]
	ErrMongoInsertError          = errors.New("[mongo] MongoDB Insert Function Occured Error")
	ErrMongoFindError            = errors.New("[mongo] MongoDB Find Function Occured Error")
	ErrMongoUpdateError          = errors.New("[mongo] MongoDB Update Function Occured Error")
	ErrMongoDeleteError          = errors.New("[mongo] MongoDB Delete Function Occured Error")
	ErrMongoInvalidObjectIDError = errors.New("[mongo] MongoDB ObjectID is invalid")

	// [users] & [logins]
	ErrUserNotRegistered     = errors.New("[users] User Not Registered Before")
	ErrUserAlreadyRegistered = errors.New("[users] User Already Registered Before")
	ErrUserPasswordIncorrect = errors.New("[users] User Password Incorrect")

	// [tokens]
	ErrTokenNotRegistered = errors.New("[tokens] Token Not Registered Before or Expired")
)
