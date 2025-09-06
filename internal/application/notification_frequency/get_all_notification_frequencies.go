package notification_frequency

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetAllNotificationFrequenciesHandler struct {
	storage  db.NotificationFrequencyRepository
	observer *observability.Observability
}

func NewGetAllNotificationFrequenciesHandler(
	storage db.NotificationFrequencyRepository,
	observer *observability.Observability,
) *GetAllNotificationFrequenciesHandler {
	return &GetAllNotificationFrequenciesHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *GetAllNotificationFrequenciesHandler) Handle(
	ctx context.Context) ([]*app_model.ApplicationNotificationFrequency, error) {
	frequencies, err := h.storage.GetAllNotificationFrequencies(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get notification frequencies")
		return nil, fmt.Errorf("failed to get notification frequencies: %w", err)
	}

	return app_model.NewApplicationNotificationFrequencies(frequencies), nil
}
