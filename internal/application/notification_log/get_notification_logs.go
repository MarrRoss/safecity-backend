package notification_log

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"time"
)

type GetNotificationLogsHandler struct {
	notificationLogsStorage db.NotificationLogRepository
}

func NewGetNotificationLogsHandler(
	notificationLogsStorage db.NotificationLogRepository) *GetNotificationLogsHandler {
	return &GetNotificationLogsHandler{
		notificationLogsStorage: notificationLogsStorage,
	}
}

func (h *GetNotificationLogsHandler) Handle(ctx context.Context) ([]*app_model.ApplicationNotificationLogType, error) {
	notifLogs, err := h.notificationLogsStorage.GetAllLogs(ctx)
	if err != nil {
		return nil, err
	}
	var notifLogResponses []*app_model.ApplicationNotificationLogType
	for _, notifLog := range notifLogs {
		notifLogResponses = append(notifLogResponses, &app_model.ApplicationNotificationLogType{
			SendTime: notifLog.CreatedAt.Format(time.RFC3339),
			Context:  notifLog.Context,
		})
	}
	return notifLogResponses, nil
}
