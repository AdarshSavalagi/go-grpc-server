package db

import (
	"log"
	"pg-consumer/internal/config"
	migrations_orm "pg-consumer/internal/orm/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(config *config.TDatabaseConfig) (*gorm.DB, error) {

	dsn := config.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if true {
		log.Println("🔄 Performing database migrations...")
		if err := migrations_orm.Migrate(db); err != nil {
			log.Printf("❌ Failed to perform database migrations: %v", err)
			return nil, err
		}
		log.Println("✅ Database migrations completed successfully")
	}else{
		log.Println("Database migration is set to false")
	}

	return db, nil
}
