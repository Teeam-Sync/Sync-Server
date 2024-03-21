package grpc_handler

import (
	"context"
	"slices"
	"strings"

	"github.com/Teeam-Sync/Sync-Server/api/converter"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	jwtService "github.com/Teeam-Sync/Sync-Server/server/service/jwt"
	utils_errors "github.com/Teeam-Sync/Sync-Server/utils/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	// Interceptor에서 token 검증을 하지 않아도 되는 요청
	passingInterceptor = []string{
		"/proto.v1.AuthService/SignUp",
		"/proto.v1.AuthService/SignIn",
	}
)

func authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if slices.Contains(passingInterceptor, info.FullMethod) { // passing interceptor
		logger.Debug(info.FullMethod, " : Passing authUnaryInterceptor")

		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok { // missing metadata
		logger.Debug("Metadata Missing: ", info.FullMethod)
		return nil, utils_errors.GrpcErrInvalidMetadata.GetGrpcError(nil)
	}

	var authHeader string

	if md["authorization"] != nil {
		authHeader = md["authorization"][0]
	} else if md["Authorization"] != nil {
		authHeader = md["Authorization"][0]
	} else {
		logger.Debug("authorization field not found from Metadata: ", md)
		return nil, utils_errors.GrpcErrInvalidMetadata.GetGrpcError(nil)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		logger.Debug("Bearer not found in authorization token: ", authHeader)
		return nil, utils_errors.GrpcErrInvalidMetadata.GetGrpcError(nil)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if info.FullMethod == "/proto.v1.AuthService/RefreshToken" {
		newCtx := context.WithValue(ctx, converter.RefreshTokenKey, token)
		return handler(newCtx, req)
	}
	uid, err := jwtService.VerifyAccessToken(token)
	if err == utils_errors.ErrExpiredAccessToken {
		logger.Debug("Expired Token : ", err)
		return resp, utils_errors.GrpcErrExpiredAccessToken.GetGrpcError(nil)
	} else if err != nil { // Invalid Token
		logger.Debug("Invalid Token : ", err)
		return resp, utils_errors.GrpcErrUnauthenticated.GetGrpcError(nil)
	}

	newCtx := context.WithValue(ctx, converter.UidKey, uid)

	return handler(newCtx, req)
}
