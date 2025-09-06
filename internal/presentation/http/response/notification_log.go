package response

import "awesomeProjectDDD/internal/application/app_model"

type GetNotificationLogsResponse struct {
	Logs []NotificationLog
}

type NotificationLog struct {
	ID                    string
	NotificationSettingID string
	SendTime              string
	Context               string
}

func NewGetNotificationLogsResponse(
	notificationLogs []*app_model.ApplicationNotificationLogType) *GetNotificationLogsResponse {
	var logs []NotificationLog
	for _, notifLog := range notificationLogs {
		logs = append(logs, NotificationLog{
			SendTime: notifLog.SendTime,
			Context:  notifLog.Context,
		})
	}
	return &GetNotificationLogsResponse{
		Logs: logs,
	}
}
