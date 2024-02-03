package auth

import (
	"context"
	"time"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	loginsColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/logins"
	usersColl "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/users"
	mongo_utils "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb/utils"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"github.com/Teeam-Sync/Sync-Server/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthServer struct {
	v1.UnimplementedAuthServiceServer
}

func (*AuthServer) SignUp(ctx context.Context, req *v1.SignUpRequest) (*v1.SignUpResponse, error) {
	logger.Debug(req)

	hashedPassword := utils.MakeHash(req.Password)

	err := mongo_utils.SignUp(loginsColl.LoginsSchema{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}, usersColl.UsersSchema{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err == converter.ErrUserAlreadyRegistered { // 이미 회원가입 된 유저입니다.
		return &v1.SignUpResponse{
			Status: &v1.Status{
				Success: false,
				Code:    v1.StatusCode_USER_ALREADY_REGISTERED,
			},
		}, nil
	} else if err != nil {
		return &v1.SignUpResponse{
			Status: &v1.Status{
				Success: false,
				Code:    v1.StatusCode_UNEXPECTED_ERROR,
			},
		}, nil
	}

	return &v1.SignUpResponse{
		Status: &v1.Status{
			Success: true,
			Code:    v1.StatusCode_NO_ERROR,
		},
	}, nil
}
