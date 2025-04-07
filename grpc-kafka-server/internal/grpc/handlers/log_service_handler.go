package handlers

import (
	"context"
	"log"
	"logs-grpc-server/grpc-kafka-server/internal/grpc/proto"
	kafka_util "logs-grpc-server/grpc-kafka-server/internal/kafka_util" // Alias kafka_util

	googleProto "google.golang.org/protobuf/proto" // Alias for the `proto` package
)

// LogServiceHandler implements the LogServiceServer interface
// It embeds UnimplementedLogServiceServer to satisfy the interface
type LogServiceHandler struct {
	proto.UnimplementedLogServiceServer // Embedding the unimplemented service
}

// UploadLog handles the RPC call for uploading a log event
func (h *LogServiceHandler) UploadLog(ctx context.Context, req *proto.LogEvent) (*proto.LogResponse, error) {
	// Process the log event (you can also push this to Kafka or a database)
	log.Printf("Received log event: %+v", req)

	// Serialize the Protobuf object into binary data using the aliased googleProto package
	data, err := googleProto.Marshal(req) // Using the alias for proto.Marshal
	if err != nil {
		log.Printf("Failed to serialize LogEvent: %v", err)
		return nil, err // Returning error if serialization fails
	}

	// Send the serialized binary data to Kafka
	kafka_util.SendToKafka(ctx, data)

	// Respond back with a success message
	return &proto.LogResponse{
		Success: true,
		Message: "Log uploaded successfully",
	}, nil
}

func (h *LogServiceHandler) UploadLogs(ctx context.Context, req *proto.BulkLogRequest)(*proto.LogResponse, error) {
	// Process the log event (you can also push this to Kafka or a database)
	log.Printf("Received log event: %+v", req)

	// Serialize the Protobuf object into binary data using the aliased googleProto package
	data, err := googleProto.Marshal(req) // Using the alias for proto.Marshal
	if err != nil {
		log.Printf("Failed to serialize LogEvent: %v", err)
		return nil, err // Returning error if serialization fails
	}

	// Send the serialized binary data to Kafka
	kafka_util.SendToKafka(ctx, data)

	// Respond back with a success message
	return &proto.LogResponse{
		Success: true,
		Message: "Log uploaded successfully",
	}, nil
}