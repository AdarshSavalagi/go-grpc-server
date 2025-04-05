package server

import (
	"log"
	"logs-grpc-server/grpc-kafka-server/internal/grpc/handlers"
	"logs-grpc-server/grpc-kafka-server/internal/grpc/proto"
	"logs-grpc-server/grpc-kafka-server/internal/kafka_util"
	"net"

	"google.golang.org/grpc"
)

func StartGRPCServer() {
	// Set up a listener on the desired port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}
	kafka_util.Init()

	// Create a new gRPC server instance
	srv := grpc.NewServer()

	// Register the LogService handler with the gRPC server
	proto.RegisterLogServiceServer(srv, &handlers.LogServiceHandler{})
	
	// Start the gRPC server
	log.Println("gRPC server running on port 50051...")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
