package app

import (
	"grpc-producer/internal/config"
	"grpc-producer/internal/server"
	"grpc-producer/internal/utils"
	"net"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	Config      *config.Config
	Logger      *logrus.Logger
	Server      *grpc.Server
	KafkaWriter *sarama.AsyncProducer
}

func InitApp() (*App, error) {
	logger := logrus.New()

	config, err := config.InitConfig()
	if err != nil {
		logger.Fatalf("Error initializing config: %v", err)
		return nil, err
	}
	server, err := server.InitGRPCServer(config)
	if err != nil {
		logger.Fatalf("Error initializing grpc server: %v", err)
		return nil, err
	}
	kafkaWriter, err := utils.InitKafkaWriters(&config.Kafka)
	if err != nil {
		logger.Fatalf("Error initializing kafka writer: %v", err)
		return nil, err
	}
	app := &App{
		Config:      config,
		Logger:      logger,
		Server:      server,
		KafkaWriter: kafkaWriter,
	}

	SetupHandler(app)

	return app, nil
}

func (a *App) Run() {
	port := a.Config.Server.Port
	a.Logger.Infof("Starting gRPC server on port %s...", strconv.Itoa(port))
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		a.Logger.Fatalf("Failed to listen on port %s: %v", strconv.Itoa(port), err)
		return
	}
	a.Logger.Infof("gRPC server running on port %s...", strconv.Itoa(port))
	if err := a.Server.Serve(listener); err != nil {
		a.Logger.Fatalf("Failed to serve: %v", err)
		return
	}
}
