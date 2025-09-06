package db

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type NotificationLogRepository interface {
	AddNotificationLog(ctx context.Context, notificationLog *entity.NotificationLog) error
	//GetAllLogs(ctx context.Context) ([]*entity.NotificationLog, error)
	//GetUserLogs(ctx context.Context, notifySettingID value_object.ID) ([]*entity.NotificationLog, error)
	GetLastLogByNotificationID(ctx context.Context, notifySettingID value_object.ID) (*entity.NotificationLog, error)
	GetLastLogsByNotificationIDs(ctx context.Context, ids []value_object.ID) ([]*response.NotificationLogDB, error)
}
