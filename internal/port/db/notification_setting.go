package db

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type NotificationSettingRepository interface {
	AddNotificationSetting(ctx context.Context, setting *aggregate.NotificationSetting) error
	AddZoneNotificationSetting(ctx context.Context, id, notificationID, zoneID value_object.ID) error
	AddBatteryNotificationSetting(ctx context.Context, id, notificationID value_object.ID, battery int) error
	NotificationSettingExists(ctx context.Context, setting *aggregate.NotificationSetting) (bool, error)
	GetNotificationSetting(ctx context.Context, id value_object.ID) (*response.NotificationSettingDB, error)
	GetNotificationSettingsByZoneForUser(ctx context.Context, receiverID, zoneID value_object.ID) ([]*response.NotificationSettingMinBatteryDB, error)
	//GetNotificationSettingIDByReceiverSender(ctx context.Context,
	//	receiverID, senderID value_object.externalID) (value_object.externalID, error)
	GetDangerZoneNotificationSettings(ctx context.Context, senderID value_object.ID, zoneID value_object.ID) ([]value_object.ID, error)
	GetNotificationSettingsByIDs(ctx context.Context, ids []value_object.ID) ([]*response.NotificationSettingWithFrequencyDB, error)
	GetNotificationsSettingsIDsByReceiver(ctx context.Context,
		receiverID value_object.ID) ([]value_object.ID, error)
	GetIntegratedReceiversBySettingIDs(ctx context.Context, settingIDs []value_object.ID) ([]*response.UserIntegration, error)
	//GetNotificationSettingSenderID(ctx context.Context, id value_object.ID) (value_object.ID, error)
	GetNotificationSettingsByZoneReceiverAndFamily(ctx context.Context,
		userID, familyID, zoneID value_object.ID) ([]*response.NotificationSettingDB, error)
	GetBatteryNotificationSettingsByChild(ctx context.Context, childID value_object.ID) ([]*response.BatterySettingWithFrequencyDB, error)
	GetBatteryNotificationSettingsByReceiver(ctx context.Context, parentID value_object.ID) ([]*response.NotificationSettingMinBatteryDB, error)
	FindLocationSenderIDsByReceiver(ctx context.Context, receiverID value_object.ID) ([]value_object.ID, error)
	UpdateNotificationSetting(ctx context.Context, setting *aggregate.NotificationSetting) error
	//DeleteNotificationSetting(ctx context.Context, id value_object.externalID) error
}
