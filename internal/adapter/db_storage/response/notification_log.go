package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type NotificationLogDB struct {
	ID             string    `db:"id"`
	NotificationID string    `db:"notification_id"`
	CreatedAt      time.Time `db:"created_at"`
	Context        string    `db:"context"`
}

func NotificationLogDbToEntity(db *NotificationLogDB) (*entity.NotificationLog, error) {
	id, err := value_object.NewIDFromString(db.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert log id: %w", err)
	}
	notifyID, err := value_object.NewIDFromString(db.NotificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert notify id: %w", err)
	}
	return &entity.NotificationLog{
		ID:                    id,
		NotificationSettingID: notifyID,
		CreatedAt:             db.CreatedAt,
		Context:               db.Context,
	}, nil
}
