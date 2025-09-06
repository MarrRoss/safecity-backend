package response

type NotificationSettingMessengerTypeDB struct {
	NotifySettingID string `db:"id_notify_setting"`
	MessengerTypeID string `db:"id_mes_type"`
	MessengerType   string `db:"mes_type"`
}
