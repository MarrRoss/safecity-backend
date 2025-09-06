package request

import "github.com/google/uuid"

type GetNotificationSettingRequest struct {
	UserID   string `params:"user_id"`
	FamilyID string `params:"family_id"`
	ZoneID   string `params:"zone_id"`
}

type GetNotificationSettingsRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name GetNotificationSettingsRequest

type GetNotificationsSendersByReceiverRequest struct {
	UserID string `json:"receiver_id" path:"receiver_id"`
}

type GetAvailableNotificationSendersRequest struct {
	FamilyID uuid.UUID `params:"id" format:"uuid"`
} // @name GetAvailableNotificationSendersRequest

type GetAvailableZonesRequest struct {
	ReceiverID string `json:"receiver_id" path:"receiver_id"`
	SenderID   string `json:"sender_id" path:"sender_id"`
}

type NotificationEventType string // @name NotificationEventType

const (
	ZoneNotificationEventType    NotificationEventType = "zone"
	BatteryNotificationEventType NotificationEventType = "battery"
)

type AddNotificationSettingRequest struct {
	FrequencyID *uuid.UUID            `json:"frequency_id,omitempty" format:"uuid"`
	EventType   NotificationEventType `json:"event_type" enum:"zone,battery"`
	SenderID    uuid.UUID             `json:"sender_id" format:"uuid"`
	ZoneID      *uuid.UUID            `json:"zone_id,omitempty" format:"uuid"`
	Battery     *int                  `json:"min_battery,omitempty" format:"int"`
} // @name AddNotificationSettingRequest

type ChangeNotificationSettingRequest struct {
	ID          string    `json:"id"`
	FrequencyID *string   `json:"frequency_id"`
	NotTypesIDs *[]string `json:"notify_types_ids"`
	MesTypesIDs *[]string `json:"mes_types_ids"`
}

type DeleteNotificationSettingRequest struct {
	ID string `params:"id"`
}
