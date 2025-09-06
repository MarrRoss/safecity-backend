package zone

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"github.com/google/uuid"
)

type DeleteZoneQuery struct {
	ID     uuid.UUID
	UserID string
}

type DeleteZoneHandler struct {
	zoneStorage db.ZoneRepository
	userStorage db.UserRepository
	observer    *observability.Observability
}

func NewDeleteZoneHandler(
	zoneStorage db.ZoneRepository,
	userStorage db.UserRepository,
	observer *observability.Observability,
) *DeleteZoneHandler {
	return &DeleteZoneHandler{
		zoneStorage: zoneStorage,
		userStorage: userStorage,
		observer:    observer,
	}
}

func (h *DeleteZoneHandler) Handle(ctx context.Context, query DeleteZoneQuery) error {
	extID, err := value_object.NewIDFromString(query.UserID)
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

	id, err := value_object.NewIDFromString(query.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
		return err
	}
	zone, err := h.zoneStorage.GetZone(ctx, id)
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
	err = zone.StopExistence()
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to stop zone existence")
		return err
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
