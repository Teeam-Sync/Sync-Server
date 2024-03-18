package auth

import (
	"context"
	"time"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	loginsColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/users"
	authService "github.com/Teeam-Sync/Sync-Server/server/service/auth"
	"github.com/Teeam-Sync/Sync-Server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthServer struct {
	v1.UnimplementedAuthServiceServer
}

// [error] GrpcErrUserAlreadyRegistered, GrpcErrUnexpected
func (*AuthServer) SignUp(ctx context.Context, req *v1.SignUpRequest) (res *v1.SignUpResponse, err error) {
	logger.Debug("[SignUp] started : ", req)

	hashedPassword := utils.MakeHash(req.Password)

	err = authService.SignUp(loginsColl.LoginSchema{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}, usersColl.UserSchema{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err == converter.ErrUserAlreadyRegistered { // 이미 회원가입 된 유저입니다.
		return nil, converter.GrpcErrUserAlreadyRegistered.GetGrpcError(nil)
	} else if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at authService.SingUp"
		return nil, converter.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	return &v1.SignUpResponse{
		Token: "",
	}, nil
}

// [error] GrpcErrUserNotRegistered, GrpcErrUserPasswordIncorrect, GrpcErrUnexpected
func (*AuthServer) SignIn(ctx context.Context, req *v1.SignInRequest) (*v1.SignInResponse, error) {
	logger.Debug("[SignUp] started : ", req)

	hashedPassword := utils.MakeHash(req.Password)

	err := authService.SignIn(loginsColl.LoginSchema{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err == converter.ErrUserNotRegistered { // 회원가입되지 않은 유저인 경우
		return nil, converter.GrpcErrUserNotRegistered.GetGrpcError(nil)
	} else if err == converter.ErrUserPasswordIncorrect { // 비밀번호가 올바르지 않은 경우
		return nil, converter.GrpcErrUserPasswordIncorrect.GetGrpcError(nil)
	} else if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at authService.SignIn"
		return nil, converter.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	return &v1.SignInResponse{
		Token: "",
	}, nil
}
