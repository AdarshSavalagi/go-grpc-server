package migrations_orm

import (
	context_orm "pg-consumer/internal/orm/models/contexts"
	events_orm "pg-consumer/internal/orm/models/events"
	raw_logs_orm "pg-consumer/internal/orm/models/raw_logs"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {

	return db.AutoMigrate(
		&context_orm.Context{},
		&events_orm.Event{},
		&raw_logs_orm.RawLog{},
	)
}
