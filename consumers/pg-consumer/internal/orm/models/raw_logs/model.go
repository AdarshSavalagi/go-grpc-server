package raw_logs_orm

import (
	"time"

	"gorm.io/datatypes"
)

// Enum definition in Go matching Protobuf LogLevel
type LogLevel int32

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	CRASH
)

type RawLog struct {
	UUID          string         `gorm:"primaryKey;type:uuid"`
	ContextID     string         `gorm:"column:context_id;type:varchar(100);"` // FK to context table
	Level         LogLevel       `gorm:"column:level;type:int;"`
	Message       string         `gorm:"column:message;type:text"`
	StackTrace    string         `gorm:"column:stack_trace;type:text"`
	LoggedAt      time.Time      `gorm:"column:logged_at;type:timestamp;"`
	LogProperties datatypes.JSON `gorm:"column:log_properties;type:jsonb;"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
