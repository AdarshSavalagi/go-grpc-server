package contexts_orm

import (
	"time"
)

type Context struct {
	ContextID        string   `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"column:user_id;type:varchar(100);"`
	SessionID string `gorm:"column:session_id;type:varchar(100);"`

	// Device Info
	DeviceID      string `gorm:"column:device_id;type:varchar(100)"`
	DeviceType    string `gorm:"column:device_type;type:varchar(100)"`
	DeviceFamily  string `gorm:"column:device_family;type:varchar(100)"`
	DeviceCarrier string `gorm:"column:device_carrier;type:varchar(100)"`

	// OS & App Info
	Platform     string `gorm:"column:platform;type:varchar(50)"`
	OSName       string `gorm:"column:os_name;type:varchar(50)"`
	OSVersion    string `gorm:"column:os_version;type:varchar(50)"`
	AppID        int64  `gorm:"column:app_id;type:bigint"`
	AppVersion   string `gorm:"column:app_version;type:varchar(50)"`
	StartVersion string `gorm:"column:start_version;type:varchar(50)"`
	SDKLibrary   string `gorm:"column:sdk_library;type:varchar(100)"`

	// Location & Locale
	Language    string  `gorm:"column:language;type:varchar(20)"`
	IPAddress   string  `gorm:"column:ip_address;type:varchar(45)"`
	City        string  `gorm:"column:city;type:varchar(100)"`
	Region      string  `gorm:"column:region;type:varchar(100)"`
	Country     string  `gorm:"column:country;type:varchar(100)"`
	LocationLat float64 `gorm:"column:location_lat"`
	LocationLng float64 `gorm:"column:location_lng"`

	// Event Timestamp
	EventTime time.Time `gorm:"column:event_time;type:timestamp;"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
