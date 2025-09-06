package response

import "time"

//type NotificationSettingFrequencySenderDB struct {
//	externalID          string `db:"notification_setting_id"`
//	FrequencyID string `db:"frequency_id"`
//	SenderID    string `db:"sender_id"`
//}

type NotificationSettingDB struct {
	ID          string    `db:"id"`
	ReceiverID  string    `db:"receiver_id"`
	SenderID    string    `db:"sender_id"`
	EventType   string    `db:"event_type"`
	FrequencyID *string   `db:"frequency_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
type NotificationSettingMinBatteryDB struct {
	ID          string    `db:"id"`
	ReceiverID  string    `db:"receiver_id"`
	SenderID    string    `db:"sender_id"`
	EventType   string    `db:"event_type"`
	FrequencyID *string   `db:"frequency_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	MinBattery  *int      `db:"battery_threshold"`
}

type NotificationSettingWithFrequencyDB struct {
	ID         string         `db:"id"`
	ReceiverID string         `db:"receiver_id"`
	SenderID   string         `db:"sender_id"`
	EventType  string         `db:"event_type"`
	Frequency  *time.Duration `db:"frequency"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}

type BatterySettingWithFrequencyDB struct {
	NotifySettingID string        `db:"notification_id"`
	Threshold       int           `db:"battery_threshold"`
	Frequency       time.Duration `db:"frequency"`
}

type BatterySettingWithFrequencySenderDB struct {
	NotifySettingID string        `db:"notification_id"`
	Threshold       int           `db:"battery_threshold"`
	FrequencyID     string        `db:"frequency_id"`
	Frequency       time.Duration `db:"frequency"`
	SenderID        string        `db:"sender_id"`
}

type NotificationSettingsIDsDB struct {
	ID string `db:"id"`
}
