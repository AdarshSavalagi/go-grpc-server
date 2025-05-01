package utils

import (
	"errors"
	"grpc-producer/internal/config"
	"log"

	"github.com/IBM/sarama"
)

func InitKafkaWriters(cfg *config.TKafkaConfig) (*sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = cfg.RetryAttempts
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.ChannelBufferSize = cfg.BufferChannelSize
	config.Producer.MaxMessageBytes = cfg.BufferChannelSize
	if len(cfg.Topics) == 0 {
		return nil, errors.New("no topics configured")
	}

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-producer.Successes():
				log.Printf("✅ Sent to Topic: %q partition %d at offset %d data : %s \n", msg.Topic,msg.Partition, msg.Offset, msg.Value)
			case err := <-producer.Errors():
				log.Printf("❌ Kafka error: %v\n", err)
			}
		}
	}()

	return &producer, nil
}
