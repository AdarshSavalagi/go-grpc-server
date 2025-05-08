package events_orm

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

// Insert inserts a single Event into the database.
func Insert(db *gorm.DB, event *Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}
	if event.ContextID == "" {
		return errors.New("event ID cannot be empty")
	}
	log.Printf("Inserting event %v", event)
	return db.Create(event).Error
}

// BulkInsert inserts multiple events in a single batch.
func BulkInsert(db *gorm.DB, events []Event) error {
	if len(events) == 0 {
		return errors.New("empty event list")
	}
	return db.CreateInBatches(events, 100).Error
}

// GetByID fetches an event by its UUID.
func GetByID(db *gorm.DB, id string) (*Event, error) {
	var event Event
	if err := db.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetByContextID fetches all events for a given context ID.
func GetByContextID(db *gorm.DB, contextID string) ([]Event, error) {
	var events []Event
	err := db.Where("context_id = ?", contextID).Find(&events).Error
	return events, err
}

// GetByEventType fetches events filtered by event type.
func GetByEventType(db *gorm.DB, eventType string) ([]Event, error) {
	var events []Event
	err := db.Where("event_type = ?", eventType).Find(&events).Error
	return events, err
}

// UpdateMessage updates the message and display name of a given event.
func UpdateMessage(db *gorm.DB, id string, newMessage, newDisplayName string) error {
	return db.Model(&Event{}).Where("id = ?", id).Updates(map[string]interface{}{
		"message":      newMessage,
		"display_name": newDisplayName,
	}).Error
}

// DeleteByID deletes an event by its ID.
func DeleteByID(db *gorm.DB, id string) error {
	return db.Delete(&Event{}, "id = ?", id).Error
}

// Exists checks if an event with a given ID exists.
func Exists(db *gorm.DB, id string) (bool, error) {
	var count int64
	if err := db.Model(&Event{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// WithTx allows running actions within a DB transaction context.
func WithTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}