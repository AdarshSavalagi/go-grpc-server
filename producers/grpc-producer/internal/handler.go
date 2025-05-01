package app

import (
	"grpc-producer/internal/handlers"
	"grpc-producer/internal/proto/logs-grpc/gen/go/models"
)

func SetupHandler(app *App) {
	models.RegisterLoggingServiceServer(app.Server, &handlers.GrpcHandler{
		KafkaWriters: app.KafkaWriter,
		Logger:       app.Logger,
		Config:       app.Config,
	})
}
