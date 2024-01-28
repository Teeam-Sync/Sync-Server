package grpc_handler

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Teeam-Sync/Sync-Server/api/handler/health_check"
	"github.com/Teeam-Sync/Sync-Server/api/proto/v1/v1connect"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	lis = flag.String("lis", "50051", "listen address")
)

func Initialize() error {
	port := os.Getenv("PORT")
	if port == "" {
		logger.Debug("port is empty")
		port = *lis
	}

	healthCheck := &health_check.HealthCheckServer{}
	mux := http.NewServeMux()
	path, handler := v1connect.NewHealthCheckServiceHandler(healthCheck)
	mux.Handle(path, handler)
	logger.Info("connectRPC Server is listening on ", port)
	http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

	return nil
}
