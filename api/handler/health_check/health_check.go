package health_check

import (
	"context"

	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
)

type HealthCheckServer struct {
	v1.UnimplementedHealthCheckServiceServer
}

func (*HealthCheckServer) Check(ctx context.Context, req *v1.CheckRequest) (*v1.CheckResponse, error) {
	logger.Debug(req.Hi)

	return &v1.CheckResponse{
		Bye: "bye",
	}, nil
}
