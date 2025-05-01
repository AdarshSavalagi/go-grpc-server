package server

import (
	"grpc-producer/internal/config"

	"google.golang.org/grpc"
)

func InitGRPCServer(config *config.Config) (*grpc.Server, error) {
	// Initialize the gRPC server
	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(config.Server.MaxRecSize),
		grpc.MaxSendMsgSize(config.Server.MaxSendSize),
	)
	return srv, nil

}
