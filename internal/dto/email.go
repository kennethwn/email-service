package dto

import (
	"time"

	"gorm.io/gorm"
)

type EmailHistoryStatus uint

const (
	EmailHistoryPending EmailHistoryStatus = 0
	EmailHistorySuccess EmailHistoryStatus = 1
	EmailHistoryFailed  EmailHistoryStatus = 2
)

var EmailHistoryStatusToString = map[EmailHistoryStatus]string{
	EmailHistoryPending: "PENDING",
	EmailHistorySuccess: "SUCCESS",
	EmailHistoryFailed:  "FAILED",
}

var EmailHistoryStatusTypeSelector = map[string]EmailHistoryStatus{
	"PENDING": EmailHistoryPending,
	"SUCCESS": EmailHistorySuccess,
	"FAILED":  EmailHistoryFailed,
}

type EmailHistory struct {
	ID        string         `gorm:"id,primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"created_at,index" json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at,index" json:"deleted_at"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Subject   string         `json:"subject"`
	Body      string         `json:"body"`
	Status    uint           `json:"status"`
	IsActive  bool           `json:"is_active"`
}

type EmailTask struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type RetryEmailRequest struct {
	ID string `json:"id"`
}
