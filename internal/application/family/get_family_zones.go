package family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"github.com/google/uuid"
)

type GetFamilyZonesQuery struct {
	ID uuid.UUID
}

type GetFamilyZonesHandler struct {
	familyStorage db.FamilyRepository
	observer      *observability.Observability
}

func NewGetFamilyZonesHandler(
	familyStorage db.FamilyRepository,
	observer *observability.Observability,
) *GetFamilyZonesHandler {
	return &GetFamilyZonesHandler{
		familyStorage: familyStorage,
		observer:      observer,
	}
}

func (h *GetFamilyZonesHandler) Handle(ctx context.Context,
	query GetFamilyZonesQuery) ([]*app_model.ApplicationZone, error) {
	familyID, err := value_object.NewIDFromString(query.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return nil, err
	}
	family, err := h.familyStorage.GetFamily(ctx, familyID)
	if err != nil {
		if errors.Is(err, adapter.ErrFamilyNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("family not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
		return nil, err
	}

	zones, err := h.familyStorage.GetFamilyZones(ctx, family.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrZoneNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("zones not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family zones")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family zones")
		return nil, err
	}
	return app_model.NewApplicationZones(zones), nil
}
