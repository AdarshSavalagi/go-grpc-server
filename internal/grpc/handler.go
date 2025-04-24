package grpc

import (
	context_service "logs-grpc-server/internal/grpc/handlers/context"
	event_service "logs-grpc-server/internal/grpc/handlers/event"
	log_service "logs-grpc-server/internal/grpc/handlers/log"
	log "logs-grpc-server/internal/grpc/proto"
)

func SetupGrpcHandlers(app *App) {
	// Register the LogServiceHandler with the gRPC server

	log.RegisterContextServiceServer(app.Server, &context_service.ContextServiceHandler{
		KafkaWriters: app.KafkaWriter,
		Logger:       app.Logger,
	})

	log.RegisterEventServiceServer(app.Server, &event_service.EventServiceHandler{
		KafkaWriters: app.KafkaWriter,
		Logger:       app.Logger,
	})

	log.RegisterLogServiceServer(app.Server, &log_service.LogServiceHandler{
		KafkaWriters: app.KafkaWriter,
		Logger:       app.Logger,
	})
}
