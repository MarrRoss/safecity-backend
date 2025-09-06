package request

type GetUserNotificationLogsRequest struct {
	UserID string `json:"user_id" path:"user_id"`
}

type AddNotificationLogRequest struct {
	NotificationSettingID string `json:"notification_setting_id" path:"notification_setting_id"`
	SendTime              string `json:"send_time" path:"send_time"`
	Context               string `json:"context" path:"context"`
}
