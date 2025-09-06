package messenger_type

import (
	"awesomeProjectDDD/internal/observability"
	"context"
)

type GetAllMessengerTypesHandler struct {
	//storage  db.MessengerTypeRepository
	observer *observability.Observability
}

func NewGetAllMessengerTypesHandler(
	//storage db.MessengerTypeRepository,
	observer *observability.Observability) *GetAllMessengerTypesHandler {
	return &GetAllMessengerTypesHandler{
		//storage:  storage,
		observer: observer,
	}
}

func (h *GetAllMessengerTypesHandler) Handle(
	ctx context.Context) error {
	//messengerTypes, err := h.storage.GetAllMessengerTypes(ctx)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get messenger types")
	//	return nil, fmt.Errorf("failed to get messenger types: %w", err)
	//}
	//
	//return app_model.NewApplicationMessengerTypes(messengerTypes), nil
	return nil
}
