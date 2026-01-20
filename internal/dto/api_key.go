package dto

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ApiKey struct {
	ID           string         `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	KeyHash      string         `json:"key_hash"`
	Name         string         `json:"name"`
	AllowedIPs   pq.StringArray `gorm:"type:text[]" json:"allowed_ips"`
	MaxPerMinute int            `json:"max_per_minute"`
	IsActive     bool           `json:"is_active"`
}

type VerifyAPIKeyResponse struct {
	IsValid            bool
	ServiceName        string
	ThresholdRateLimit int
	APIKey             string
}
