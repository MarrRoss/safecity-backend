package family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AddFamilyZoneCommand struct {
	FamilyID   uuid.UUID
	Name       string
	Boundaries []app_model.ApplicationLocation
	Safety     bool
}

type AddFamilyZoneHandler struct {
	familyStorage db.FamilyRepository
	zoneStorage   db.ZoneRepository
	observer      *observability.Observability
}

func NewAddFamilyZoneHandler(
	familyStorage db.FamilyRepository,
	zoneStorage db.ZoneRepository,
	observer *observability.Observability,
) *AddFamilyZoneHandler {
	return &AddFamilyZoneHandler{
		familyStorage: familyStorage,
		zoneStorage:   zoneStorage,
		observer:      observer,
	}
}

func (h *AddFamilyZoneHandler) Handle(ctx context.Context, cmd AddFamilyZoneCommand) (value_object.ID, error) {
	familyID, err := value_object.NewIDFromString(cmd.FamilyID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return value_object.ID{}, err
	}
	family, err := h.familyStorage.GetFamily(ctx, familyID)
	if err != nil {
		if errors.Is(err, adapter.ErrFamilyNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("family not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
		return value_object.ID{}, err
	}

	var boundaries []value_object.Location
	for _, loc := range cmd.Boundaries {
		lat := value_object.Latitude(loc.Latitude)
		lng := value_object.Longitude(loc.Longitude)
		location := value_object.Location{
			Latitude:  lat,
			Longitude: lng,
		}
		boundaries = append(boundaries, location)
	}

	overlaps, err := h.zoneStorage.ZoneOverlaps(ctx, family.ID, &boundaries)
	if overlaps {
		h.observer.Logger.Trace().Err(err).Msg("zone overlaps other zones")
		return value_object.ID{}, application.ErrZoneOverlapsZones
	}

	name, err := value_object.NewZoneName(cmd.Name)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add zone name")
		return value_object.ID{}, fmt.Errorf("failed to add zone name: %w", err)
	}
	safety, err := value_object.NewSafety(cmd.Safety)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add zone safety")
		return value_object.ID{}, fmt.Errorf("failed to add zone safety: %w", err)
	}
	zone, err := entity.NewZone(name, &boundaries, safety, family.ID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create zone")
		return value_object.ID{}, fmt.Errorf("failed to create zone: %w", err)
	}

	err = h.zoneStorage.AddZone(ctx, zone)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding zone")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding zone")
		return value_object.ID{}, err
	}
	return zone.ID, nil
}
