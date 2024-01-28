package health_check

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
)

type HealthCheckServer struct {
}

func (*HealthCheckServer) Check(ctx context.Context, req *connect.Request[v1.CheckRequest]) (*connect.Response[v1.CheckResponse], error) {
	logger.Debug("Request headers: ", req.Header())
	res := connect.NewResponse(&v1.CheckResponse{
		Bye: "bye",
	})
	res.Header().Set("SyncServer-Version", "v1")
	return res, nil
}
