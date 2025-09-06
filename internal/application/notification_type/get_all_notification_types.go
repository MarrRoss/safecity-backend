package notification_type

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetAllNotificationTypesHandler struct {
	storage  db.NotificationTypeRepository
	observer *observability.Observability
}

func NewGetAllNotificationTypesHandler(
	storage db.NotificationTypeRepository,
	observer *observability.Observability) *GetAllNotificationTypesHandler {
	return &GetAllNotificationTypesHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *GetAllNotificationTypesHandler) Handle(
	ctx context.Context) ([]*app_model.ApplicationNotificationType, error) {
	notificationTypes, err := h.storage.GetAllNotificationTypes(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get notification types")
		return nil, fmt.Errorf("failed to get notification types: %w", err)
	}

	return app_model.NewApplicationNotificationTypes(notificationTypes), nil
}
