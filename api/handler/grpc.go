package grpc_handler

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/Teeam-Sync/Sync-Server/api/handler/auth"
	"github.com/Teeam-Sync/Sync-Server/api/handler/health_check"
	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	v1.RegisterHealthCheckServiceServer(s, &health_check.HealthCheckServer{})
	v1.RegisterAuthServiceServer(s, &auth.AuthServer{})
	reflection.Register(s)
	logger.Debug("server listening at ", lis.Addr())

	if err := s.Serve(lis); err != nil {
		logger.Debug("server stopped with error: %v", err) // Change this line
	} else {
		logger.Debug("server stopped gracefully")
	}

	return err
}
