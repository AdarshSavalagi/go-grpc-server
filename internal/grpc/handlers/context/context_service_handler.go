package context_service

import (
	"context"

	log "logs-grpc-server/internal/grpc/proto"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

type ContextServiceHandler struct {
	KafkaWriters *sarama.AsyncProducer
	Logger       *logrus.Logger
	log.UnimplementedContextServiceServer
}

// Handle sending a single Context
func (h *ContextServiceHandler) SendContext(ctx context.Context, req *log.Context) (*log.Response, error) {
	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal Context: %v", err)
		return &log.Response{Success: false, Message: "Failed to encode context"}, err
	}

	select {
	case (*h.KafkaWriters).Input() <- &sarama.ProducerMessage{
		Topic: "contexts",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Infof("Context pushed to Kafka")
		return &log.Response{Success: true, Message: "Context sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full, context dropped")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}

// Handle sending a list of Contexts
func (h *ContextServiceHandler) SendContextList(ctx context.Context, req *log.ContextList) (*log.Response, error) {
	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal ContextList: %v", err)
		return &log.Response{Success: false, Message: "Failed to encode context list"}, err
	}

	select {
	case (*h.KafkaWriters).Input() <- &sarama.ProducerMessage{
		Topic: "contexts",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Infof("Context list pushed to Kafka")
		return &log.Response{Success: true, Message: "Context list sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full, context list dropped")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}

// Handle sending a ContextFile (raw binary content)
func (h *ContextServiceHandler) SendContextFile(ctx context.Context, req *log.ContextFileRequest) (*log.Response, error) {
	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal ContextFileRequest: %v", err)
		return &log.Response{Success: false, Message: "Failed to encode context file"}, err
	}

	select {
	case (*h.KafkaWriters).Input() <- &sarama.ProducerMessage{
		Topic: "context-files",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Infof("Context file pushed to Kafka: %s", req.FileName)
		return &log.Response{Success: true, Message: "Context file sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full, context file dropped")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}