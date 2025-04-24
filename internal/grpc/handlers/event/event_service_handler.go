package event_service

import (
	"context"
	log "logs-grpc-server/internal/grpc/proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type EventServiceHandler struct {
	KafkaWriters map[string]*kafka.Producer
	Logger       *logrus.Logger
	log.UnimplementedEventServiceServer
}

func (h *EventServiceHandler) SendEvent(ctx context.Context, req *log.EventMessage) (*log.Response, error) {
	// producer, ok := h.KafkaWriters["event"]
	// if !ok {
	// 	h.Logger.Error("Kafka producer for 'event' not found")
	// 	return &log.Response{Success: false, Message: "Kafka producer not available"}, nil
	// }

	// topic := "events" // Update this if the topic name is dynamic or from config

	// // Delivery report handler for produced messages
	// deliveryChan := make(chan kafka.Event)

	// err := producer.Produce(&kafka.Message{
	// 	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 	Key:            []byte(req.Key),
	// 	Value:          []byte(req.Value),
	// }, deliveryChan)

	// if err != nil {
	// 	h.Logger.Errorf("Failed to produce message: %v", err)
	// 	return &log.Response{Success: false, Message: "Failed to produce message"}, err
	// }

	// select {
	// case e := <-deliveryChan:
	// 	m := e.(*kafka.Message)
	// 	if m.TopicPartition.Error != nil {
	// 		h.Logger.Errorf("Delivery failed: %v", m.TopicPartition.Error)
	// 		return &log.Response{Success: false, Message: "Delivery failed"}, m.TopicPartition.Error
	// 	}
	// 	h.Logger.Infof("Message delivered to topic %s [%d] at offset %v",
	// 		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	// case <-ctx.Done():
	// 	h.Logger.Warn("Context cancelled before message delivery")
	// 	return &log.Response{Success: false, Message: "Context cancelled"}, ctx.Err()
	// case <-time.After(10 * time.Second):
	// 	h.Logger.Warn("Timeout waiting for delivery report")
	// 	return &log.Response{Success: false, Message: "Timeout waiting for delivery report"}, nil
	// }

	// close(deliveryChan)
	return &log.Response{Success: true, Message: "Event sent successfully"}, nil
}
