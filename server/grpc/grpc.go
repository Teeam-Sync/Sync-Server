package grpc_handler

import (
	"flag"
	"fmt"
	"net"
	"os"

	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
	logger "github.com/Teeam-Sync/Sync-Server/logging"
	"github.com/Teeam-Sync/Sync-Server/server/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	lis = flag.String("lis", "50051", "listen address")
)

func Initialize() error {
	port := os.Getenv("PORT")
	if port == "" {
		logger.Info("port is empty")
		port = *lis
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(authUnaryInterceptor))
	v1.RegisterAuthServiceServer(s, &auth.AuthServer{})
	reflection.Register(s)
	logger.Info("server listening at ", lis.Addr())

	if err := s.Serve(lis); err != nil {
		logger.Info("server stopped with error: %v", err) // Change this line
	} else {
		logger.Info("server stopped gracefully")
	}

	return err
}
