package log_service

import (
	"context"
	log "logs-grpc-server/internal/grpc/proto"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type LogServiceHandler struct {
	KafkaWriters *sarama.AsyncProducer
	Logger       *logrus.Logger
	log.UnimplementedLogServiceServer
}


func (h *LogServiceHandler)SendLogList(ctx context.Context, req *log.LogList)(*log.Response,error){
	return &log.Response{Success: true, Message: "Log list sent successfully"}, nil
}
func (h *LogServiceHandler) SendLog(ctx context.Context, req *log.LogMessage) (*log.Response, error) {
	// producer, ok := h.KafkaWriters["log"]
	// if !ok {
	// 	h.Logger.Error("Kafka producer for 'log' not found")
	// 	return &log.Response{Success: false, Message: "Kafka producer not available"}, nil
	// }

	// topic := "logs" // Change this if you need a configurable topic

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
	// 	h.Logger.Warn("Context cancelled before delivery")
	// 	return &log.Response{Success: false, Message: "Context cancelled"}, ctx.Err()
	// case <-time.After(10 * time.Second):
	// 	h.Logger.Warn("Timeout waiting for delivery report")
	// 	return &log.Response{Success: false, Message: "Timeout waiting for delivery report"}, nil
	// }

	// close(deliveryChan)
	return &log.Response{Success: true, Message: "Log sent successfully"}, nil
}
