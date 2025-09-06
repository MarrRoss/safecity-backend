package notification_log

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetUserNotificationLogsQuery struct {
	UserID string
}

type GetUserNotificationLogsHandler struct {
	notificationLogsStorage     db.NotificationLogRepository
	notificationSettingsStorage db.NotificationSettingRepository
	userStorage                 db.UserRepository
}

func NewGetUserNotificationLogsHandler(
	notificationLogsStorage db.NotificationLogRepository,
	notificationSettingsStorage db.NotificationSettingRepository) *GetUserNotificationLogsHandler {
	return &GetUserNotificationLogsHandler{
		notificationLogsStorage:     notificationLogsStorage,
		notificationSettingsStorage: notificationSettingsStorage,
	}
}

func (h *GetUserNotificationLogsHandler) Handle(ctx context.Context,
	query GetUserNotificationLogsQuery) ([]*app_model.ApplicationNotificationLogType, error) {
	receiverID, err := value_object.NewIDFromString(query.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	_, err = h.userStorage.GetUser(ctx, receiverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	//notifySettingsIDs, err := h.notificationSettingsStorage.GetNotificationsSettingsIDsByReceiver(ctx, u.ID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get notification settings IDs: %w", err)
	//}
	//var allNotifyLogs []*entity.NotificationLog
	//for _, notifySettingID := range notifySettingsIDs {
	//	settingNotifyLogs, err := h.notificationLogsStorage.GetUserLogs(ctx, notifySettingID)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to get notification log: %w", err)
	//	}
	//	allNotifyLogs = append(allNotifyLogs, settingNotifyLogs...)
	//}

	var logsModels [][]string
	//for _, notifyLog := range allNotifyLogs {
	//	senderID, err := h.notificationSettingsStorage.GetNotificationSettingSenderID(
	//		ctx, *notifyLog.NotificationSettingID)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to get notification setting: %w", err)
	//	}
	//	//sender, err := h.userStorage.GetUser(ctx, senderID)
	//	//if err != nil {
	//	//	return nil, fmt.Errorf("failed to get user: %w", err)
	//	//}
	//	//senderName := sender.Name.FirstName.String() + " " + sender.Name.LastName.String() +
	//	//	" " + sender.Name.Patronymic.String()
	//	//senderPhone := sender.Phone.String()
	//}
	return app_model.NewApplicationNotificationLogTypes(logsModels), nil
}
