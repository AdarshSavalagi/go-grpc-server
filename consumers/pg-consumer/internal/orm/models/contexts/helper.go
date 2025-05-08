package contexts_orm

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

// InsertContext inserts a new context record
func InsertContext(db *gorm.DB, ctx *Context) error {
	if ctx == nil {
		return errors.New("context cannot be nil")
	}
	log.Printf("Inserting context %v", ctx)
	return db.Create(ctx).Error
}

// GetContextByID fetches a context by primary key
func GetContextByID(db *gorm.DB, id uint) (*Context, error) {
	var ctx Context
	err := db.First(&ctx, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ctx, err
}

// GetContextBySessionID fetches a context by session ID
func GetContextBySessionID(db *gorm.DB, sessionID string) (*Context, error) {
	var ctx Context
	err := db.Where("session_id = ?", sessionID).First(&ctx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ctx, err
}

// UpdateContext updates a context record
func UpdateContext(db *gorm.DB, ctx *Context) error {
	return db.Save(ctx).Error
}

// DeleteContext deletes a context record by ID
func DeleteContext(db *gorm.DB, id uint) error {
	return db.Delete(&Context{}, id).Error
}