package app

import (
	"context"
	"log"
	"pg-consumer/internal/config"
	"pg-consumer/internal/db"
	"pg-consumer/internal/kafka_util"
	"sync"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type App struct {
	KafkaConsumer sarama.Consumer
	DB            *gorm.DB
	Config        *config.Config
}

func InitializeApp() (*App, error) {
	config, err := config.Init()
	if err != nil {
		return nil, err
	}

	db, err := db.InitPostgres(&config.Database)
	if err != nil {
		return nil, err
	}

	consumer, err := kafka_util.InitKafkaConsumer(&config.Kafka)
	if err != nil {
		return nil, err
	}

	return &App{
		KafkaConsumer: consumer,
		DB:            db,
		Config:        config,
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("🚀 Starting consumers...")
	kafka_util.ConsumeMessages(ctx, a.KafkaConsumer, a.DB, a.Config, wg)
}
