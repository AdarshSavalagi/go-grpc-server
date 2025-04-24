package grpc

import (
	"logs-grpc-server/internal/grpc/config"
	"logs-grpc-server/internal/grpc/server"
	"logs-grpc-server/internal/kafka_util"
	"net"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	Config      *config.Config
	Logger      *logrus.Logger
	Server      *grpc.Server
	KafkaWriter map[string]*kafka.Producer
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
		logger.Fatalf("Error initializing gRPC server: %v", err)
		return nil, err
	}
	kafkaWriter, err := kafka_util.InitKafkaWriters(&config.Kafka)
	if err != nil {
		logger.Fatalf("Error initializing Kafka writers: %v", err)
		return nil, err
	}

	app := &App{
		Config:      config,
		Logger:      logger,
		Server:      server,
		KafkaWriter: kafkaWriter,
	}

	SetupGrpcHandlers(app)
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
