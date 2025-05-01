package kafka_util

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"pg-consumer/internal/config"
	contexts_orm "pg-consumer/internal/orm/models/contexts"
	events_orm "pg-consumer/internal/orm/models/events"
	raw_logs_orm "pg-consumer/internal/orm/models/raw_logs"
	"sync"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

func InitKafkaConsumer(cfg *config.TKafkaConfig) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	if len(cfg.Brokers) == 0 {
		return nil, errors.New("no brokers configured")
	}

	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama consumer: %v", err)
		return nil, err
	}

	return consumer, nil
}

func ConsumeMessages(ctx context.Context, consumer sarama.Consumer, db *gorm.DB, config *config.Config, wg *sync.WaitGroup) {
	for _, topic := range config.Kafka.Topics {
		wg.Add(1)
		go func(topicName string) {
			defer wg.Done()
			log.Printf("🔄 Starting consumer for topic: %s\n", topicName)

			partitions, err := consumer.Partitions(topicName)
			if err != nil {
				log.Printf("❌ Failed to get partitions for topic %s: %v", topicName, err)
				return
			}

			for _, partition := range partitions {
				pc, err := consumer.ConsumePartition(topicName, partition, sarama.OffsetNewest)
				if err != nil {
					log.Printf("❌ Failed to consume partition %d for topic %s: %v", partition, topicName, err)
					continue
				}

				wg.Add(1)
				go func(pc sarama.PartitionConsumer, partition int32) {
					defer wg.Done()
					defer pc.Close()

					log.Printf("✅ Listening on topic %s partition %d\n", topicName, partition)

					for {
						select {
						case msg := <-pc.Messages():
							log.Printf("📨 [%s] Msg: %s", topicName, string(msg.Value))
							switch topicName {
							case "logs":
								var rawLog raw_logs_orm.RawLog
								if err := json.Unmarshal(msg.Value, &rawLog); err == nil {
									err = raw_logs_orm.Insert(db, &rawLog)
								}

							case "events":
								var event events_orm.Event
								if err := json.Unmarshal(msg.Value, &event); err == nil {
									err = events_orm.Insert(db, &event)
								}

							case "contexts":
								var context contexts_orm.Context
								if err := json.Unmarshal(msg.Value, &context); err == nil {
									err = contexts_orm.InsertContext(db, &context)
								}
							case "logs-batch":
								var batch struct {
									Logs []raw_logs_orm.RawLog `json:"logs"`
								}
								if err := json.Unmarshal(msg.Value, &batch); err == nil {
									for _, logEntry := range batch.Logs {
										if insertErr := raw_logs_orm.Insert(db, &logEntry); insertErr != nil {
											log.Printf("❌ Insert error for log UUID %s: %v", logEntry.UUID, insertErr)
										}
									}
								}
							case "contexts-batch":
								var batch struct {
									Contexts []contexts_orm.Context `json:"contexts"`
								}
								if err := json.Unmarshal(msg.Value, &batch); err == nil {
									for _, context := range batch.Contexts {
										if insertErr := contexts_orm.InsertContext(db, &context); insertErr != nil {
											log.Printf("❌ Insert error for log UUID %s: %v", context.ContextID, insertErr)
										}
									}
								}

							case "events-batch":
								var batch struct {
									Contexts []events_orm.Event `json:"events"`
								}
								if err := json.Unmarshal(msg.Value, &batch); err == nil {
									for _, event := range batch.Contexts {
										if insertErr := events_orm.Insert(db, &event); insertErr != nil {
											log.Printf("❌ Insert error for log UUID %s: %v", event.UUID, insertErr)
										}
									}
								}


							}

							if err != nil {
								log.Printf("❌ DB Insert Error on topic %s: %v", topicName, err)
							}
						case err := <-pc.Errors():
							log.Printf("❌ Kafka error on topic %s: %v", topicName, err)
						case <-ctx.Done():
							log.Printf("🛑 Stopping consumer for topic %s", topicName)
							return
						}
					}
				}(pc, partition)
			}
		}(topic)
	}
}
