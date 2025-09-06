package db

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type NotificationFrequencyRepository interface {
	GetFrequency(ctx context.Context, id value_object.ID) (*entity.NotificationFrequency, error)
	GetAllNotificationFrequencies(ctx context.Context) ([]*entity.NotificationFrequency, error)
	GetNotificationFrequenciesByIDs(ctx context.Context, ids []string) ([]*entity.NotificationFrequency, error)
}
