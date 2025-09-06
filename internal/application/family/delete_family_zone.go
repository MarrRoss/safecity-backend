package family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"github.com/google/uuid"
)

type DeleteFamilyZoneCommand struct {
	FamilyID uuid.UUID
	ZoneID   uuid.UUID
}

type DeleteFamilyZoneHandler struct {
	familyStorage db.FamilyRepository
	zoneStorage   db.ZoneRepository
	observer      *observability.Observability
}

func NewDeleteFamilyZoneHandler(
	familyStorage db.FamilyRepository,
	zoneStorage db.ZoneRepository,
	observer *observability.Observability,
) *DeleteFamilyZoneHandler {
	return &DeleteFamilyZoneHandler{
		familyStorage: familyStorage,
		zoneStorage:   zoneStorage,
		observer:      observer,
	}
}

func (h *DeleteFamilyZoneHandler) Handle(ctx context.Context, cmd DeleteFamilyZoneCommand) error {
	familyID, err := value_object.NewIDFromString(cmd.FamilyID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return err
	}
	_, err = h.familyStorage.GetFamily(ctx, familyID)
	if err != nil {
		if errors.Is(err, adapter.ErrFamilyNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("family not found")
			return err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family")
			return err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
		return err
	}
	zoneID, err := value_object.NewIDFromString(cmd.ZoneID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
		return err
	}
	_, err = h.zoneStorage.GetZone(ctx, zoneID)
	if err != nil {
		if errors.Is(err, adapter.ErrZoneNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("zone not found")
			return err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting zone")
			return err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone")
		return err
	}

	/// заменить на update
	//err = h.familyZoneStorage.DeleteFamilyZone(ctx, family.externalID, zone.externalID)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while delete zone from family")
	//		return err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while delete zone from family")
	//	return err
	//}
	return nil
}
