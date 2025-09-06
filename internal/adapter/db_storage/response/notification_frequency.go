package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type NotificationFrequencyDB struct {
	Id        uuid.UUID     `db:"id"`
	Frequency time.Duration `db:"frequency"`
}

func NotificationFrequencyDbToEntity(db *NotificationFrequencyDB) (*entity.NotificationFrequency, error) {
	id, err := value_object.NewIDFromString(db.Id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to frequency id")
	}
	frequency := db.Frequency
	return &entity.NotificationFrequency{
		ID:        id,
		Frequency: frequency,
	}, nil
}

func NotificationFrequencyDbListToEntityList(
	dbList []*NotificationFrequencyDB) ([]*entity.NotificationFrequency, error) {
	var frequencies []*entity.NotificationFrequency
	for _, db := range dbList {
		frequency, err := NotificationFrequencyDbToEntity(db)
		if err != nil {
			return nil, fmt.Errorf("failed to convert storage data to entity: %v", err)
		}
		frequencies = append(frequencies, frequency)
	}
	return frequencies, nil
}
