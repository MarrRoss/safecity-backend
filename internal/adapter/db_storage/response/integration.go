package response

import (
	"github.com/google/uuid"
	"time"
)

type SystemIntegration struct {
	ID         int    `db:"id"`
	ExternalID string `db:"external_id"`
}

type UserIntegration struct {
	UserID          uuid.UUID         `db:"user_id"`
	NotifySettingID uuid.UUID         `db:"notify_setting_id"`
	System          SystemIntegration `db:"system"`
}

type IntegrationDB struct {
	ID         int       `db:"id"`
	UserID     string    `db:"user_id"`
	SystemID   string    `db:"system_id"`
	ExternalID *string   `db:"external_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
