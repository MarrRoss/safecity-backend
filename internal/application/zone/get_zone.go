package zone

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetZoneQuery struct {
	ID string
}

type GetZoneHandler struct {
	storage db.ZoneRepository
}

func NewGetZoneHandler(storage db.ZoneRepository) *GetZoneHandler {
	return &GetZoneHandler{
		storage: storage,
	}
}

func (h *GetZoneHandler) Handle(ctx context.Context, query GetZoneQuery) (*app_model.ApplicationZone, error) {
	id, err := value_object.NewIDFromString(query.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid zone id: %w", err)
	}
	zone, err := h.storage.GetZone(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone: %w", err)
	}
	return app_model.NewApplicationZone(zone), nil
}
