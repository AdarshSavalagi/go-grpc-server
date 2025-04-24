package server

import (
	"logs-grpc-server/internal/grpc/config"

	"google.golang.org/grpc"
)

func InitGRPCServer(config *config.Config) (*grpc.Server, error) {
	// Initialize the gRPC server
	srv := grpc.NewServer()
	return srv, nil

}
