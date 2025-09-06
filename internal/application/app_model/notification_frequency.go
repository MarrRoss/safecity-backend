package app_model

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"github.com/google/uuid"
	"time"
)

type ApplicationNotificationFrequency struct {
	ID        uuid.UUID
	Frequency time.Duration
}

func NewApplicationNotificationFrequency(frequency *entity.NotificationFrequency) *ApplicationNotificationFrequency {
	return &ApplicationNotificationFrequency{
		ID:        frequency.ID.ToRaw(),
		Frequency: frequency.Frequency,
	}
}

func NewApplicationNotificationFrequencies(
	frequencies []*entity.NotificationFrequency) []*ApplicationNotificationFrequency {
	appFrequencies := make([]*ApplicationNotificationFrequency, 0)
	for _, frequency := range frequencies {
		appFrequencies = append(appFrequencies, NewApplicationNotificationFrequency(frequency))
	}
	return appFrequencies
}
