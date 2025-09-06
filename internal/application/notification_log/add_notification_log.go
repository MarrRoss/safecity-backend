package notification_log

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
	"time"
)

type AddNotificationLogCommand struct {
	NotificationSettingID string
	SendTime              string
	Context               string
}

type AddNotificationLogHandler struct {
	notificationLogsStorage     db.NotificationLogRepository
	notificationSettingsStorage db.NotificationSettingRepository
}

func NewAddNotificationLogHandler(
	notificationLogsStorage db.NotificationLogRepository,
	notificationSettingsStorage db.NotificationSettingRepository) *AddNotificationLogHandler {
	return &AddNotificationLogHandler{
		notificationLogsStorage: notificationLogsStorage,
	}
}

func (h *AddNotificationLogHandler) Handle(ctx context.Context, cmd AddNotificationLogCommand) error {
	notifSetID, err := value_object.NewIDFromString(cmd.NotificationSettingID)
	if err != nil {
		return fmt.Errorf("invalid notification setting id: %w", err)
	}
	_, err = h.notificationSettingsStorage.GetNotificationSetting(ctx, notifSetID)
	if err != nil {
		return fmt.Errorf("failed to get notification setting: %w", err)
	}
	sendTime, err := time.Parse(time.RFC3339, cmd.SendTime)
	if err != nil {
		return fmt.Errorf("failed to parse send time: %w", err)
	}
	notifLog, err := entity.NewNotificationLog(
		&notifSetID,
		sendTime,
		cmd.Context,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification log: %w", err)
	}
	err = h.notificationLogsStorage.AddNotificationLog(ctx, notifLog)
	if err != nil {
		return fmt.Errorf("failed to add notification log: %w", err)
	}
	return nil
}
