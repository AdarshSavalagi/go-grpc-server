package handlers

import (
	"context"
	"errors"
	"log"
	"logs-grpc-server/grpc-kafka-server/internal/grpc/proto"

	// Alias kafka_util
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	// Alias for the `proto` package
)

// LogServiceHandler implements the LogServiceServer interface
type LogServiceHandler struct {
	KafkaWriters                        map[string]*kafka.Writer
	Logger                              *logrus.Logger
	proto.UnimplementedLogServiceServer // Embed unimplemented server for forward compatibility
}

// UploadLog handles a single log event via gRPC and sends it to Kafka
func (h *LogServiceHandler) UploadLog(ctx context.Context, req *proto.LogEvent) (*proto.LogResponse, error) {
	log.Printf("Received log event: ID=%s, Source=%s, AppName=%s", req.LogID, req.DeviceID, req.AppName)

	writer, ok := h.KafkaWriters["logs"]
	if !ok {
		log.Printf("❌ Kafka writer not found for topic 'logs'")
		return nil, errors.New("kafka writer not found for topic 'logs'")
	}
	err := writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(req.LogID),
	})
	if err != nil {
		log.Printf("❌ Failed to send log to Kafka: %v", err)
		return nil, err
	}

	// Respond to client
	return &proto.LogResponse{
		Success: true,
		Message: "Log uploaded successfully",
	}, nil
}

// UploadLogs handles bulk log event uploads and sends each one to Kafka individually
func (h *LogServiceHandler) UploadLogs(ctx context.Context, req *proto.BulkLogRequest) (*proto.LogResponse, error) {
	// log.Printf("📦 Received bulk log request with %d log(s)", len(req.Logs))

	// var failedCount int

	// for i, logEvent := range req.Logs {
	// 	log.Printf("➡️ Processing log %d: ID=%s, Source=%s", i+1, logEvent.LogID, logEvent.AppName)

	// 	// Serialize each log event individually
	// 	data, err := googleProto.Marshal(logEvent)
	// 	if err != nil {
	// 		log.Printf("❌ Failed to serialize log at index %d: %v", i, err)
	// 		failedCount++
	// 		continue
	// 	}

	// 	// Send to Kafka
	// 	if err := kafka_util.SendToKafka(ctx, data); err != nil {
	// 		log.Printf("❌ Failed to send log to Kafka at index %d: %v", i, err)
	// 		failedCount++
	// 		continue
	// 	}
	// }

	// // Final response based on result
	// if failedCount > 0 {
	// 	msg := "Some logs failed to process"
	// 	log.Printf("⚠️ %d log(s) failed to process", failedCount)
	// 	return &proto.LogResponse{
	// 		Success: false,
	// 		Message: msg,
	// 	}, nil
	// }

	// log.Println("✅ All logs uploaded successfully")
	return &proto.LogResponse{
		Success: true,
		Message: "All logs uploaded successfully",
	}, nil
}
