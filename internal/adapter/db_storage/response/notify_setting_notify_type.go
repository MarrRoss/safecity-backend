package response

type NotificationSettingNotificationTypeDB struct {
	NotifySettingID string `db:"id_notify_setting"`
	NotifyTypeID    string `db:"id_notify_type"`
	NotifyType      string `db:"notify_type"`
}
