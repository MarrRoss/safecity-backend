package zone

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ChangeZoneCommand struct {
	ID         uuid.UUID
	Name       *string
	Boundaries *[]app_model.ApplicationLocation
	Safety     *bool
	UserID     string
}

type ChangeZoneHandler struct {
	zoneStorage db.ZoneRepository
	userStorage db.UserRepository
	observer    *observability.Observability
}

func NewChangeZoneHandler(
	zoneStorage db.ZoneRepository,
	userStorage db.UserRepository,
	observer *observability.Observability,
) *ChangeZoneHandler {
	return &ChangeZoneHandler{
		zoneStorage: zoneStorage,
		userStorage: userStorage,
		observer:    observer,
	}
}

func (h *ChangeZoneHandler) Handle(ctx context.Context, cmd ChangeZoneCommand) error {
	extID, err := value_object.NewIDFromString(cmd.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return err
	}
	_, err = h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return err
	}
	zoneID, err := value_object.NewIDFromString(cmd.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
		return err
	}
	zone, err := h.zoneStorage.GetZone(ctx, zoneID)
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

	if cmd.Name != nil {
		name, err := value_object.NewZoneName(*cmd.Name)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to add zone name")
			return fmt.Errorf("failed to add zone name: %w", err)
		}
		err = zone.ChangeName(name)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to change zone name")
			return fmt.Errorf("failed to change zone name: %w", err)
		}
	}
	if cmd.Safety != nil {
		safety, err := value_object.NewSafety(*cmd.Safety)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to add zone safety")
			return fmt.Errorf("failed to add zone safety: %w", err)
		}
		err = zone.ChangeSafety(safety)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to change zone safety")
			return fmt.Errorf("failed to change zone safety: %w", err)
		}
	}

	if cmd.Boundaries != nil {
		var boundaries []value_object.Location
		for _, loc := range *cmd.Boundaries {
			lat := value_object.Latitude(loc.Latitude)
			lng := value_object.Longitude(loc.Longitude)
			location := value_object.Location{
				Latitude:  lat,
				Longitude: lng,
			}
			boundaries = append(boundaries, location)
		}

		err = zone.ChangeBoundaries(&boundaries)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to change zone boundaries")
			return fmt.Errorf("failed to change zone boundaries: %w", err)
		}
	}
	err = h.zoneStorage.UpdateZone(ctx, zone)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while updating zone")
			return err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while updating zone")
		return err
	}
	return nil
}
