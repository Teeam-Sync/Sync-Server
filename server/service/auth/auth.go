package authService

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"github.com/Teeam-Sync/Sync-Server/server/database/mongodb"
	loginsColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/users"
	jwtService "github.com/Teeam-Sync/Sync-Server/server/service/jwt"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// [logins & users] logins & users 컬렉션에 유저를 가입시키는 함수 (transcation 포함)
/* ErrUnexpectedError, ErrUserAlreadyRegistered, ErrMongoInsertError */
func SignUp(loginUser loginsColl.LoginSchema, user usersColl.UserSchema) (token converter.JWTToken, err error) {
	loginRepository := loginsColl.NewMongoLoginRepository()
	userRepository := usersColl.NewMongoUserRepository()

	uid := primitive.NewObjectID()
	loginUser.Uid = uid
	user.Uid = uid

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := mongodb.Client.StartSession()
	if err != nil { // unexpected error
		logger.Error(err)
		return token, utils_errors.ErrUnexpectedError
	}

	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		err = loginRepository.InsertLoginUser(ctx, loginUser)
		if err != nil { // 이미 가입되었거나 unexpected error
			return nil, err
		}

		err = userRepository.InsertUser(ctx, user)
		if err != nil { // unexpected error
			logger.Error(err)
			return nil, utils_errors.ErrUnexpectedError
		}

		return nil, nil
	}, txnOptions)
	if err == utils_errors.ErrUserAlreadyRegistered { // 이미 가입되어있는 경우
		return token, utils_errors.ErrUserAlreadyRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return token, utils_errors.ErrMongoInsertError
	}

	token, err = jwtService.CreateJWTToken(uid.Hex())
	if err != nil { // unexpected error
		logger.Error(err)
		return token, utils_errors.ErrUnexpectedError
	}

	return token, nil
}

// [logins] logins 컬렉션에서 유저를 확인하고 로그인하는 함수
/* ErrUnexpectedError, ErrUserNotRegistered, ErrMongoFindError, ErrUserPasswordIncorrect */
func SignIn(loginUser loginsColl.LoginSchema) (token converter.JWTToken, err error) {
	loginRepository := loginsColl.NewMongoLoginRepository()

	user, err := loginRepository.FindLoginUserByEmail(loginUser.Email)
	if err == utils_errors.ErrUserNotRegistered { // login Collection에 등록되어있지 않은 유저인 경우
		return token, utils_errors.ErrUserNotRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return token, utils_errors.ErrMongoFindError
	}

	if user.Password != loginUser.Password { // 비밀번호가 틀렸을 때
		return token, utils_errors.ErrUserPasswordIncorrect
	}

	token, err = jwtService.CreateJWTToken(user.Uid.Hex())
	if err != nil { // unexpected error
		logger.Error(err)
		return token, utils_errors.ErrUnexpectedError
	}

	return token, nil
}
