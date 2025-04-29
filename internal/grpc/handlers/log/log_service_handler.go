package log_service

import (
	"context"
	log "logs-grpc-server/internal/grpc/proto"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

type LogServiceHandler struct {
	KafkaWriters *sarama.AsyncProducer
	Logger       *logrus.Logger
	log.UnimplementedLogServiceServer
}

// SendLogList handles a batch of log messages and sends to Kafka
func (h *LogServiceHandler) SendLogList(ctx context.Context, req *log.LogList) (*log.Response, error) {
	producer := (*h.KafkaWriters)

	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal log list to JSON: %v", err)
		return &log.Response{Success: false, Message: "Encoding error"}, err
	}

	select {
	case producer.Input() <- &sarama.ProducerMessage{
		Topic: "logs-list",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Infof("Log list sent to Kafka")
		return &log.Response{Success: true, Message: "Log list sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full — dropping log list")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}

// SendLog handles individual log message
func (h *LogServiceHandler) SendLog(ctx context.Context, req *log.LogMessage) (*log.Response, error) {
	producer := (*h.KafkaWriters)

	data, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.Errorf("Failed to marshal log message to JSON: %v", err)
		return &log.Response{Success: false, Message: "Encoding error"}, err
	}

	select {
	case producer.Input() <- &sarama.ProducerMessage{
		Topic: "logs-message",
		Value: sarama.ByteEncoder(data),
	}:
		h.Logger.Infof("Single log message sent to Kafka")
		return &log.Response{Success: true, Message: "Log sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full — dropping log")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}

// SendLogFile handles full log file uploads (binary content)
func (h *LogServiceHandler) SendLogFile(ctx context.Context, req *log.LogFileRequest) (*log.Response, error) {
	producer := (*h.KafkaWriters)

	if req.FileContent == nil {
		h.Logger.Error("Log file content is nil")
		return &log.Response{Success: false, Message: "Empty file content"}, nil
	}

	message := &sarama.ProducerMessage{
		Topic: "logs-file",
		Value: sarama.ByteEncoder(req.FileContent),
		Headers: []sarama.RecordHeader{
			{Key: []byte("sessionId"), Value: []byte(req.SessionID)},
			{Key: []byte("fileName"), Value: []byte(req.FileName)},
		},
	}

	select {
	case producer.Input() <- message:
		h.Logger.Infof("Log file '%s' sent to Kafka", req.FileName)
		return &log.Response{Success: true, Message: "Log file sent successfully"}, nil
	default:
		h.Logger.Warn("Kafka buffer full — dropping log file")
		return &log.Response{Success: false, Message: "Kafka buffer full"}, nil
	}
}