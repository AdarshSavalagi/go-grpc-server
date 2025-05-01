package events_orm

import (
	"time"

	"gorm.io/datatypes"
)

type Event struct {
	UUID              string         `gorm:"column:uuid;primaryKey;type:uuid;"`                 // Primary key UUID
	ContextID       string         `gorm:"column:context_id;type:varchar(100);not null;index"`                       // Foreign key / grouping ID
	EventType       string         `gorm:"column:event_type;type:varchar(50);"`                        // E.g., LOGIN, ERROR
	DisplayName     string         `gorm:"column:display_name;type:varchar(150);"`                           // User-friendly label
	Message         string         `gorm:"column:message;type:text;"`                                        // Main message body
	EventTime       string         `gorm:"column:timestamp;type:varchar(50);"`                               // ISO string, if needed for parsing
	EventProperties datatypes.JSON `gorm:"column:event_properties;type:jsonb;"`                              // JSONB field for extra data

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`                                                    // Auto timestamps
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`                                                    // Auto timestamps
}