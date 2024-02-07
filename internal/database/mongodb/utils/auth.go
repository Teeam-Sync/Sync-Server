package mongo_utils

import (
	"context"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	"github.com/Teeam-Sync/Sync-Server/internal/database/mongodb"
	loginsColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/users"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// [logins & users] logins & users 컬렉션에 유저를 가입시키는 함수 (transcation 포함)
func SignUp(loginUser loginsColl.LoginsSchema, user usersColl.UsersSchema) (err error) {

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := mongodb.Client.StartSession()
	if err != nil { // unexpected error
		logger.Error(err)
		return err
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		uid := primitive.NewObjectID()
		loginUser.Uid = uid
		user.Uid = uid

		err = loginsColl.InsertLoginUser(ctx, loginUser)
		if err != nil { // 이미 가입되었거나 unexpected error
			return nil, err
		}

		err = usersColl.InsertUser(ctx, user)
		if err != nil { // unexpected error
			logger.Error(err)
			return nil, converter.ErrUnexpectedError
		}

		return nil, nil
	}, txnOptions)
	if err == converter.ErrUserAlreadyRegistered { // 이미 가입되어있는 경우
		return err
	} else if err != nil { // unexpected error
		logger.Error(err)
		return converter.ErrMongoInsertError
	}

	return nil
}

// [logins] logins 컬렉션에서 유저를 확인하고 로그인하는 함수
func SignIn(loginUser loginsColl.LoginsSchema) (err error) {
	findedUser, err := loginsColl.FindLoginUserByEmail(loginUser.Email)
	if err == converter.ErrUserNotRegistered { // login Collection에 등록되어있지 않은 유저인 경우
		return err
	} else if err != nil { // unexpected error
		logger.Error(err)
		return converter.ErrUnexpectedError
	}

	if findedUser.Password != loginUser.Password { // 비밀번호가 틀렸을 때
		return converter.ErrUserPasswordIncorrect
	}
	return nil
}
