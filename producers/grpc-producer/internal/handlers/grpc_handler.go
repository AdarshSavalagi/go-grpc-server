package handlers

import (
	"context"
	"grpc-producer/internal/config"
	"grpc-producer/internal/proto/logs-grpc/gen/go/models"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type GrpcHandler struct {
	Logger       *logrus.Logger
	KafkaWriters *sarama.AsyncProducer
	Config       *config.Config
	models.UnimplementedLoggingServiceServer
}

// Helper
func (h *GrpcHandler) sendToKafka(topic string, msg proto.Message) (*models.UploadResponse, error) {
	data, err := protojson.Marshal(msg)
	if err != nil {
		h.Logger.Errorf("Failed to marshal message: %v", err)
		return nil, err
	}

	select {
	case (*h.KafkaWriters).Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}:
		return &models.UploadResponse{Success: true, Message: "Message sent", LiveTrackingEnabled: false}, nil
	default:
		h.Logger.Warn("Kafka buffer full, message dropped")
		return &models.UploadResponse{Success: false, Message: "Kafka buffer full", LiveTrackingEnabled: false}, nil
	}
}

// === Single Uploads ===

func (h *GrpcHandler) UploadLog(ctx context.Context, req *models.LogMessage) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["logs"], req)
}

func (h *GrpcHandler) UploadEvent(ctx context.Context, req *models.EventMessage) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["events"], req)
}

func (h *GrpcHandler) UploadContext(ctx context.Context, req *models.ContextMessage) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["contexts"], req)
}

// === Batch Uploads ===

func (h *GrpcHandler) UploadLogs(ctx context.Context, req *models.LogBatch) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["logs-batch"], req)
}

func (h *GrpcHandler) UploadEvents(ctx context.Context, req *models.EventBatch) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["events-batch"], req)
}

func (h *GrpcHandler) UploadContexts(ctx context.Context, req *models.ContextBatch) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["contexts-batch"], req)
}

// === File Uploads ===

func (h *GrpcHandler) UploadLogsFile(ctx context.Context, req *models.CompressedFileUpload) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["logs-file"], req)
}

func (h *GrpcHandler) UploadEventsFile(ctx context.Context, req *models.CompressedFileUpload) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["events-file"], req)
}

func (h *GrpcHandler) UploadContextsFile(ctx context.Context, req *models.CompressedFileUpload) (*models.UploadResponse, error) {
	return h.sendToKafka(h.Config.Kafka.Topics["contexts-file"], req)
}

// === Streaming Uploads ===

func (h *GrpcHandler) StreamLogs(stream models.LoggingService_StreamLogsServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		resp, err := h.UploadLog(stream.Context(), msg)
		if err != nil {
			return err
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func (h *GrpcHandler) StreamEvents(stream models.LoggingService_StreamEventsServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		resp, err := h.UploadEvent(stream.Context(), msg)
		if err != nil {
			return err
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func (h *GrpcHandler) StreamContexts(stream models.LoggingService_StreamContextsServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		resp, err := h.UploadContext(stream.Context(), msg)
		if err != nil {
			return err
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
