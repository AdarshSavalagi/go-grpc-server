package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"logs-grpc-server/grpc-kafka-server/internal/grpc/proto"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	protobuf "google.golang.org/protobuf/proto"
)

type LogServiceHandler struct {
	KafkaWriters                        map[string]*kafka.Writer
	Logger                              *logrus.Logger
	proto.UnimplementedLogServiceServer // Embed unimplemented server for forward compatibility
}

// UploadLog handles single log upload
func (h *LogServiceHandler) UploadLog(ctx context.Context, req *proto.LogEvent) (*proto.LogResponse, error) {
	// log.Printf("Received log event: ID=%s, Source=%s, AppName=%s", req.LogID, req.DeviceID, req.AppName)

	// writer, ok := h.KafkaWriters["logs"]
	// if !ok {
	// 	log.Printf("❌ Kafka writer not found for topic 'logs'")
	// 	return nil, errors.New("kafka writer not found for topic 'logs'")
	// }

	// // Serialize the log event using protobuf (binary format)
	// logData, err := protobuf.Marshal(req) // Using protobuf's Marshal function
	// if err != nil {
	// 	log.Printf("❌ Failed to serialize log event: %v", err)
	// 	return nil, err
	// }

	// // Send the serialized log to Kafka
	// err = writer.WriteMessages(ctx, kafka.Message{
	// 	Value: logData, // Send the protobuf-encoded log event
	// })
	// if err != nil {
	// 	log.Printf("❌ Failed to send log to Kafka: %v", err)
	// 	return nil, err
	// }

	// // Respond to client
	// log.Printf("✅ Log uploaded successfully: %s", req.LogID)
	return &proto.LogResponse{
		Success: true,
		Message: "Log uploaded successfully",
	}, nil
}

// UploadLogs handles bulk log uploads
func (h *LogServiceHandler) UploadLogs(ctx context.Context, req *proto.BulkLogRequest) (*proto.LogResponse, error) {
	log.Printf("Received bulk log request with %d log(s)", len(req.Logs))

	writer, ok := h.KafkaWriters["logs"]
	if !ok {
		log.Printf("❌ Kafka writer not found for topic 'logs'")
		return nil, errors.New("kafka writer not found for topic 'logs'")
	}

	var kafkaMessages []kafka.Message
	for _, logItem := range req.Logs {
		// Serialize each log event using protobuf (binary format)
		logData, err := protobuf.Marshal(logItem)
		if err != nil {
			log.Printf("❌ Failed to serialize log event: %v", err)
			continue // Skip this log if serialization fails
		}

		// Append the serialized log to the Kafka message array
		kafkaMessages = append(kafkaMessages, kafka.Message{
			Value: logData,
		})
	}

	// Write all messages to Kafka
	if len(kafkaMessages) > 0 {
		err := writer.WriteMessages(ctx, kafkaMessages...)
		if err != nil {
			log.Printf("❌ Failed to send logs to Kafka: %v", err)
			return nil, err
		}
	}

	log.Printf("✅ Successfully sent %d logs to Kafka", len(kafkaMessages))
	return &proto.LogResponse{
		Success: true,
		Message: fmt.Sprintf("%d logs uploaded successfully", len(kafkaMessages)),
	}, nil
}
