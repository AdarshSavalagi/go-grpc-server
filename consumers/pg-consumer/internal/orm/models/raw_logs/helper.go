package raw_logs_orm

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Insert inserts a single RawLog record into the DB.
func Insert(db *gorm.DB, log *RawLog) error {
	return db.Create(log).Error
}

// BulkInsert inserts multiple RawLogs efficiently in batches.
func BulkInsert(db *gorm.DB, logs []RawLog) error {
	if len(logs) == 0 {
		return errors.New("empty raw log list")
	}
	return db.CreateInBatches(logs, 100).Error
}

// GetByUUID retrieves a RawLog by its UUID.
func GetByUUID(db *gorm.DB, uuid string) (*RawLog, error) {
	var log RawLog
	err := db.First(&log, "uuid = ?", uuid).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByContextID returns all logs associated with a specific context ID.
func GetByContextID(db *gorm.DB, contextID string) ([]RawLog, error) {
	var logs []RawLog
	err := db.Where("context_id = ?", contextID).Find(&logs).Error
	return logs, err
}

// GetByLogLevel returns all logs of a specific log level.
func GetByLogLevel(db *gorm.DB, level LogLevel) ([]RawLog, error) {
	var logs []RawLog
	err := db.Where("level = ?", level).Find(&logs).Error
	return logs, err
}

// UpdateStackTrace updates only the stack trace of a specific RawLog.
func UpdateStackTrace(db *gorm.DB, uuid string, newStack string) error {
	return db.Model(&RawLog{}).
		Where("uuid = ?", uuid).
		Update("stack_trace", newStack).Error
}

// DeleteByUUID deletes a log entry by UUID.
func DeleteByUUID(db *gorm.DB, uuid string) error {
	return db.Delete(&RawLog{}, "uuid = ?", uuid).Error
}

// Exists checks if a RawLog exists by UUID.
func Exists(db *gorm.DB, uuid string) (bool, error) {
	var count int64
	err := db.Model(&RawLog{}).Where("uuid = ?", uuid).Count(&count).Error
	return count > 0, err
}

// WithTx allows executing DB logic in a transaction context.
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