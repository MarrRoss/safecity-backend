package app_model

type ApplicationNotificationLogType struct {
	SenderName string
	Phone      string
	Context    string
	Zone       string
	SendTime   string
}

func NewApplicationNotificationLogType(logModel []string) *ApplicationNotificationLogType {
	return &ApplicationNotificationLogType{
		SenderName: logModel[0],
		Phone:      logModel[1],
		Context:    logModel[2],
		Zone:       logModel[3],
		SendTime:   logModel[4],
	}
}

func NewApplicationNotificationLogTypes(logsModel [][]string) []*ApplicationNotificationLogType {
	appLogs := make([]*ApplicationNotificationLogType, len(logsModel))
	for i, log := range logsModel {
		appLogs[i] = NewApplicationNotificationLogType(log)
	}
	return appLogs
}
