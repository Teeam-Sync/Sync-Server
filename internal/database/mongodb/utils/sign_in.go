package mongo_utils

import (
	"github.com/Teeam-Sync/Sync-Server/api/converter"
	loginsColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/logins"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
)

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