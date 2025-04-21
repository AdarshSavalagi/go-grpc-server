package grpc

import (
	"logs-grpc-server/grpc-kafka-server/internal/grpc/handlers"
	"logs-grpc-server/grpc-kafka-server/internal/grpc/proto"
)

func SetupGrpcHandlers(app *App) {
	// Register the LogServiceHandler with the gRPC server
	proto.RegisterLogServiceServer(app.Server, &handlers.LogServiceHandler{
		KafkaWriters: app.KafkaWriter,
		Logger:       app.Logger,
	})
	// Register other handlers as needed
	// proto.RegisterAnotherServiceServer(app.Server, &AnotherServiceHandler{})
}
