package authService

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"github.com/Teeam-Sync/Sync-Server/server/database/mongodb"
	loginsColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// [logins & users] logins & users 컬렉션에 유저를 가입시키는 함수 (transcation 포함)
/* ErrUnexpectedError, ErrUserAlreadyRegistered, ErrMongoInsertError */
func SignUp(loginUser loginsColl.LoginSchema, user usersColl.UserSchema) (err error) {
	loginRepository := loginsColl.NewMongoLoginRepository()
	userRepository := usersColl.NewMongoUserRepository()

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := mongodb.Client.StartSession()
	if err != nil { // unexpected error
		logger.Error(err)
		return converter.ErrUnexpectedError
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		uid := primitive.NewObjectID()
		loginUser.Uid = uid
		user.Uid = uid

		err = loginRepository.InsertLoginUser(ctx, loginUser)
		if err != nil { // 이미 가입되었거나 unexpected error
			return nil, err
		}

		err = userRepository.InsertUser(ctx, user)
		if err != nil { // unexpected error
			logger.Error(err)
			return nil, converter.ErrUnexpectedError
		}

		return nil, nil
	}, txnOptions)
	if err == converter.ErrUserAlreadyRegistered { // 이미 가입되어있는 경우
		return converter.ErrUserAlreadyRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return converter.ErrMongoInsertError
	}

	return nil
}

// [logins] logins 컬렉션에서 유저를 확인하고 로그인하는 함수
/* ErrUserNotRegistered, ErrMongoFindError, ErrUserPasswordIncorrect */
func SignIn(loginUser loginsColl.LoginSchema) (err error) {
	loginRepository := loginsColl.NewMongoLoginRepository()

	findedUser, err := loginRepository.FindLoginUserByEmail(loginUser.Email)
	if err == converter.ErrUserNotRegistered { // login Collection에 등록되어있지 않은 유저인 경우
		return converter.ErrUserNotRegistered
	} else if err != nil { // unexpected error
		logger.Error(err)
		return converter.ErrMongoFindError
	}

	if findedUser.Password != loginUser.Password { // 비밀번호가 틀렸을 때
		return converter.ErrUserPasswordIncorrect
	}
	return nil
}
