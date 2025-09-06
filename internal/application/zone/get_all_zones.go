package zone

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetAllZonesHandler struct {
	storage db.ZoneRepository
}

func NewGetAllZonesHandler(storage db.ZoneRepository) *GetAllZonesHandler {
	return &GetAllZonesHandler{
		storage: storage,
	}
}

func (h *GetAllZonesHandler) Handle(ctx context.Context) ([]*app_model.ApplicationZone, error) {
	zones, err := h.storage.GetZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones: %w", err)
	}

	return app_model.NewApplicationZones(zones), nil
}
