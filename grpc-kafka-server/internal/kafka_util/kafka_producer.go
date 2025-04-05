package kafka_util

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
)

var (
	writer *kafka.Writer
)

func Init() {
	kafkaBroker := getKafkaBroker()
	topic := getKafkaTopic()

	log.Printf("Initializing Kafka writer with broker: %s, topic: %s", kafkaBroker, topic)

	// First ensure topic exists
	if err := ensureTopicExists(kafkaBroker, topic); err != nil {
		log.Printf("WARNING: Could not verify topic exists: %v", err)
	}

	writer = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaBroker),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		BatchSize:              1,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
		MaxAttempts:            3,
	}
}

func ensureTopicExists(broker, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

func SendToKafka(ctx context.Context, message []byte) error {
	log.Printf("Attempting to send message to Kafka topic: %s", writer.Topic)

	err := writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})

	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}

	log.Println("Message successfully sent to Kafka")
	return nil
}

func getKafkaBroker() string {
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		return broker
	}
	return "kafka:9092"
}

func getKafkaTopic() string {
	if topic := os.Getenv("KAFKA_TOPIC"); topic != "" {
		return topic
	}
	return "logs"
}
