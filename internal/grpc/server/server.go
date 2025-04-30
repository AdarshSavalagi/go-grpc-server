package server

import (
	"logs-grpc-server/internal/grpc/config"

	"google.golang.org/grpc"
)

func InitGRPCServer(config *config.Config) (*grpc.Server, error) {
	// Initialize the gRPC server
	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(30*1024*1024), // allow up to 16 MB
		grpc.MaxSendMsgSize(30*1024*1024),
	)
	return srv, nil

}
