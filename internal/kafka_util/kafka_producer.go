package kafka_util

import (
	"fmt"
	"log"
	"logs-grpc-server/internal/grpc/config"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ensureTopic(brokers []string, topic string, partitions, replicationFactor int) error {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers[0]})
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}
	defer admin.Close()

	topicSpec := []kafka.TopicSpecification{{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}}

	for i := 0; i < 5; i++ {
		results, err := admin.CreateTopics(nil, topicSpec, kafka.SetAdminOperationTimeout(5*time.Second))
		if err == nil && len(results) > 0 && results[0].Error.Code() == kafka.ErrNoError {
			return nil
		}
		if err != nil {
			log.Printf("Retrying topic creation for %s due to error: %v", topic, err)
		} else {
			log.Printf("Retrying topic creation for %s due to Kafka error: %v", topic, results[0].Error.String())
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to create topic %s after retries", topic)
}

func InitKafkaWriters(cfg *config.TKafkaConfig) (map[string]*kafka.Producer, error) {
	producers := make(map[string]*kafka.Producer)

	for name, topic := range cfg.Topics {
		if err := ensureTopic(cfg.Brokers, topic, 3, 1); err != nil {
			log.Printf("Warning: failed to ensure topic %s: %v", topic, err)
		}

		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers":      cfg.Brokers[0], // confluent client doesn't take []string, use comma-separated string if needed
			"acks":                   cfg.Acks,
			"queue.buffering.max.ms": cfg.WriteTimeout * 1000, // optional tuning
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create producer for topic %s: %w", topic, err)
		}

		// Optional: Delivery report handler (can help with debugging)
		go func() {
			for e := range p.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						log.Printf("Delivery failed: %v\n", ev.TopicPartition)
					}
				}
			}
		}()

		producers[name] = p
	}
	return producers, nil
}

func NewKafkaMessage(value []byte) *kafka.Message {
	return &kafka.Message{
		Value: value,
	}
}
