package zone

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"github.com/google/uuid"
)

type AddZoneCommand struct {
	Name       string
	Boundaries []app_model.ApplicationLocation
	Safety     bool
	FamilyID   uuid.UUID
}

type AddZoneHandler struct {
	zoneStorage   db.ZoneRepository
	familyStorage db.FamilyRepository
	observer      *observability.Observability
}

func NewAddZoneHandler(
	zoneStorage db.ZoneRepository,
	familyStorage db.FamilyRepository,
	observer *observability.Observability,
) *AddZoneHandler {
	return &AddZoneHandler{
		zoneStorage:   zoneStorage,
		familyStorage: familyStorage,
		observer:      observer,
	}
}

func (h *AddZoneHandler) Handle(ctx context.Context, cmd AddZoneCommand) (value_object.ID, error) {
	//id, err := value_object.NewIDFromString(cmd.FamilyID.String())
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid family id")
	//	return value_object.ID{}, err
	//}
	//family, err := h.familyStorage.GetFamily(ctx, id)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrFamilyNotFound) {
	//		h.observer.Logger.Trace().Err(err).Msg("family not found")
	//		return value_object.ID{}, err
	//	}
	//
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting family")
	//		return value_object.ID{}, err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
	//	return value_object.ID{}, err
	//}
	//name, err := value_object.NewZoneName(cmd.Name)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to add zone name")
	//	return value_object.ID{}, fmt.Errorf("failed to add zone name: %w", err)
	//}
	//safety, err := value_object.NewSafety(cmd.Safety)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to add zone safety")
	//	return value_object.ID{}, fmt.Errorf("failed to add zone safety: %w", err)
	//}
	//
	//var boundaries []value_object.Location
	//for _, loc := range cmd.Boundaries {
	//	lat := value_object.Latitude(loc.Latitude)
	//	lng := value_object.Longitude(loc.Longitude)
	//	location := value_object.Location{
	//		Latitude:  lat,
	//		Longitude: lng,
	//	}
	//	boundaries = append(boundaries, location)
	//}
	//
	//zone, err := entity.NewZone(name, &boundaries, safety, family.ID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to create zone")
	//	return value_object.ID{}, fmt.Errorf("failed to create zone: %w", err)
	//}
	//
	//err = h.zoneStorage.AddZone(ctx, zone)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while adding zone")
	//		return value_object.ID{}, err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while adding zone")
	//	return value_object.ID{}, err
	//}

	return value_object.ID{}, nil
}
