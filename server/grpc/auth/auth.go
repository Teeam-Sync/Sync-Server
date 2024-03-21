package auth

import (
	"context"
	"time"

	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	loginsColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/server/database/mongodb/users"
	authService "github.com/Teeam-Sync/Sync-Server/server/service/auth"
	jwtService "github.com/Teeam-Sync/Sync-Server/server/service/jwt"
	"github.com/Teeam-Sync/Sync-Server/utils"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServer struct {
	v1.UnimplementedAuthServiceServer
}

// GrpcErrUserAlreadyRegistered, GrpcErrUnexpected
func (*AuthServer) SignUp(ctx context.Context, req *v1.SignUpRequest) (*v1.SignUpResponse, error) {
	logger.Debug("[SignUp] started : ", req)

	hashedPassword := utils.MakeHash(req.Password)

	token, err := authService.SignUp(loginsColl.LoginSchema{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}, usersColl.UserSchema{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err == utils_errors.ErrUserAlreadyRegistered { // 이미 회원가입 된 유저입니다.
		return nil, utils_errors.GrpcErrUserAlreadyRegistered.GetGrpcError(nil)
	} else if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at authService.SingUp"
		return nil, utils_errors.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	return &v1.SignUpResponse{
		Token: token.ToPB(),
	}, nil
}

// GrpcErrUserNotRegistered, GrpcErrUserPasswordIncorrect, GrpcErrUnexpected
func (*AuthServer) SignIn(ctx context.Context, req *v1.SignInRequest) (*v1.SignInResponse, error) {
	logger.Debug("[SignIn] started : ", req)

	hashedPassword := utils.MakeHash(req.Password)

	token, err := authService.SignIn(loginsColl.LoginSchema{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err == utils_errors.ErrUserNotRegistered { // 회원가입되지 않은 유저인 경우
		return nil, utils_errors.GrpcErrUserNotRegistered.GetGrpcError(nil)
	} else if err == utils_errors.ErrUserPasswordIncorrect { // 비밀번호가 올바르지 않은 경우
		return nil, utils_errors.GrpcErrUserPasswordIncorrect.GetGrpcError(nil)
	} else if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at authService.SignIn"
		return nil, utils_errors.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	return &v1.SignInResponse{
		Token: token.ToPB(),
	}, nil
}

// GrpcErrUserNotRegistered, GrpcErrUserPasswordIncorrect, GrpcErrUnexpected
func (*AuthServer) RefreshToken(ctx context.Context, req *emptypb.Empty) (*v1.RefreshTokenResponse, error) {
	logger.Debug("[RefreshToken] started : ", req)

	refreshToken, err := utils.GetRefreshToken(ctx)
	if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at utils.GetRefreshToken"
		return nil, utils_errors.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	uid, err := jwtService.VerifyRefreshToken(refreshToken)
	if err == utils_errors.ErrExpiredRefreshToken || err == utils_errors.ErrTokenNotRegistered || err == utils_errors.ErrInvalidToken { // 만료되거나 유효하지 않은 토큰 -> 재로그인 유도
		return nil, utils_errors.GrpcErrExpiredRefreshToken.GetGrpcError(nil)
	} else if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at jwtService.VerifyRefreshToken"
		return nil, utils_errors.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	token, err := jwtService.CreateJWTToken(uid)
	if err != nil { // unexpected error
		logger.Error(err)
		errMsg := "unexpected error occured at jwtService.CreateJWTToken"
		return nil, utils_errors.GrpcErrUnexpected.GetGrpcError(&errMsg)
	}

	return &v1.RefreshTokenResponse{
		Token: token.ToPB(),
	}, nil
}
