package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"github.com/google/uuid"
)

type GetNotificationFrequencyResponse struct {
	ID        uuid.UUID `json:"id"`
	Frequency string    `json:"frequency"`
}

func NewGetNotificationFrequencyResponse(
	frequency *app_model.ApplicationNotificationFrequency) *GetNotificationFrequencyResponse {
	if frequency == nil {
		return nil
	}
	return &GetNotificationFrequencyResponse{
		ID:        frequency.ID,
		Frequency: frequency.Frequency.String(),
	}
}

func NewGetNotificationFrequenciesResponse(
	frequencies []*app_model.ApplicationNotificationFrequency) []*GetNotificationFrequencyResponse {
	notificationFrequencies := make([]*GetNotificationFrequencyResponse, 0)
	for _, frequency := range frequencies {
		notificationFrequencies = append(notificationFrequencies, NewGetNotificationFrequencyResponse(frequency))
	}
	return notificationFrequencies
}
