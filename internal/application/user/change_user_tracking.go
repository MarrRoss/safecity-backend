package user

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
)

type ChangeUserTrackingCommand struct {
	ID string
}

type ChangeUserTrackingHandler struct {
	storage  db.UserRepository
	observer *observability.Observability
}

func NewChangeUserTrackingHandler(
	storage db.UserRepository,
	observer *observability.Observability,
) *ChangeUserTrackingHandler {
	return &ChangeUserTrackingHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *ChangeUserTrackingHandler) Handle(ctx context.Context, cmd ChangeUserTrackingCommand) error {
	extID, err := value_object.NewIDFromString(cmd.ID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return err
	}
	user, err := h.storage.GetUserByExternalID(ctx, extID)
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
	err = user.ChangeTracking(!user.Tracking)
	if err != nil {
		h.observer.Logger.Trace().Msg("user is deleted")
		return err
	}

	err = h.storage.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while updating user")
			return err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while updating user")
		return err
	}

	return nil
}
