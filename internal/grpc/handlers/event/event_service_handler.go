package event_service

import (
	"context"
	log "logs-grpc-server/internal/grpc/proto"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

type EventServiceHandler struct {
	KafkaWriters *sarama.AsyncProducer
	Logger       *logrus.Logger
	log.UnimplementedEventServiceServer
}

func (h *EventServiceHandler) SendContextFile(ctx context.Context, req *log.ContextFileRequest) (*log.Response, error) {
	// Implementation for sending context file
	file := req.FileContent
	if file == nil {
		h.Logger.Error("File content is nil")
		return &log.Response{Success: false, Message: "File content is nil"}, nil
	}
	producer := (*h.KafkaWriters)

	msg := &sarama.ProducerMessage{
		Topic: "context-files",
		Key:   sarama.StringEncoder(req.FileName),
		Value: sarama.ByteEncoder(file),
		Headers: []sarama.RecordHeader{
			{Key: []byte("sessionId"), Value: []byte(req.SessionID)},
			{Key: []byte("fileName"), Value: []byte(req.FileName)},
		},
	}

	select {
	case producer.Input() <- msg:
		h.Logger.Infof("Context file queued: %s", req.FileName)
		return &log.Response{Success: true, Message: "Context file sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full, dropping context file")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}

func (h *EventServiceHandler) SendEvent(ctx context.Context, req *log.EventMessage) (*log.Response, error) {
	producer := (*h.KafkaWriters)

	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal protobuf to JSON: %v", err)
		return &log.Response{Success: false, Message: "Failed to encode message"}, err
	}

	select {
	case producer.Input() <- &sarama.ProducerMessage{
		Topic: "events",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Debugf("Event queued for Kafka topic: %s, data: %s", "events", string(data))
		return &log.Response{Success: true, Message: "Event sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka input buffer full, dropping event")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}

func (h *EventServiceHandler) SendEventList(ctx context.Context, req *log.EventList) (*log.Response, error) {
	producer := (*h.KafkaWriters)

	// Marshal the entire list of events to JSON (or your preferred format)
	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal event list to JSON: %v", err)
		return &log.Response{Success: false, Message: "Failed to encode event list"}, err
	}

	// Send the whole list of events as a single message
	select {
	case producer.Input() <- &sarama.ProducerMessage{
		Topic: "events",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Debugf("Event list queued for Kafka topic: %s, data: %s", "events", string(data))
		return &log.Response{Success: true, Message: "Event list sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka input buffer full, dropping event list")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}
