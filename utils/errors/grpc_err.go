package utils_errors

import (
	"fmt"

	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcErr struct {
	GrpcCode codes.Code
	ErrCode  v1.StatusCode
	Msg      string
}

var (
	/* 0xxxx - common */
	// Unexpected Error
	GrpcErrUnexpected = &GrpcErr{
		GrpcCode: codes.Unknown,
		ErrCode:  v1.StatusCode_UNEXPECTED_ERROR,
		Msg:      "Unexpected Error Occured",
	}
	// Required Field Missing
	GrpcErrRequiredFieldMissing = &GrpcErr{
		GrpcCode: codes.InvalidArgument,
		ErrCode:  v1.StatusCode_REQUIREDFIELD_MISSING,
		Msg:      "Required Field Missing",
	}
	// Authentication Failed
	GrpcErrUnauthenticated = &GrpcErr{
		GrpcCode: codes.Unauthenticated,
		ErrCode:  v1.StatusCode_UNAUTHENTICATED,
		Msg:      "Authentication Failed",
	}
	// Invalid Metadata(JWT Token)
	GrpcErrInvalidMetadata = &GrpcErr{
		GrpcCode: codes.Unauthenticated,
		ErrCode:  v1.StatusCode_INVALID_METADATA,
		Msg:      "Invalid Metadata(JWT Token)",
	}
	// Expired Metadata(JWT Access Token)
	GrpcErrExpiredAccessToken = &GrpcErr{
		GrpcCode: codes.Unauthenticated,
		ErrCode:  v1.StatusCode_EXPIRED_ACCESSTOKEN,
		Msg:      "Access Token Expired",
	}
	// Expired Metadata(JWT Refresh Token)
	GrpcErrExpiredRefreshToken = &GrpcErr{
		GrpcCode: codes.Unauthenticated,
		ErrCode:  v1.StatusCode_EXPIRED_REFRESHTOKEN,
		Msg:      "Refresh Token Expired",
	}

	/* 1xxxx - authService */
	// User Already Registered
	GrpcErrUserAlreadyRegistered = &GrpcErr{GrpcCode: codes.Unavailable,
		ErrCode: v1.StatusCode_USER_ALREADY_REGISTERED,
		Msg:     "User Already Registered",
	}
	// User Not Registered
	GrpcErrUserNotRegistered = &GrpcErr{GrpcCode: codes.Unavailable,
		ErrCode: v1.StatusCode_USER_NOT_REGISTERED,
		Msg:     "User Not Registered",
	}
	// User's Password Incorrect
	GrpcErrUserPasswordIncorrect = &GrpcErr{GrpcCode: codes.Unauthenticated,
		ErrCode: v1.StatusCode_USER_PASSWORD_INCORRECT,
		Msg:     "User's Password Incorrect",
	}
)

func (g *GrpcErr) GetGrpcError(msg *string) (err error) {
	if msg == nil { // 메세지가 nil로 넘어온다면, 기본 메세지 보내기
		msg = &g.Msg
	}
	err = status.Error(g.GrpcCode, fmt.Sprintf("[%d] %s", g.ErrCode, *msg))
	logger.Info(err)
	return err
}
