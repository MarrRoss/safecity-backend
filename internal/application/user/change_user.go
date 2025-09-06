package user

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ChangeUserCommand struct {
	ID        uuid.UUID
	FirstName *string
	LastName  *string
}

type ChangeUserHandler struct {
	storage  db.UserRepository
	observer *observability.Observability
}

func NewChangeUserHandler(
	storage db.UserRepository,
	observer *observability.Observability,
) *ChangeUserHandler {
	return &ChangeUserHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *ChangeUserHandler) Handle(ctx context.Context, cmd ChangeUserCommand) error {
	userID, err := value_object.NewIDFromString(cmd.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user id")
		return err
	}
	user, err := h.storage.GetUser(ctx, userID)
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
	//if user.EndedAt != nil {
	//	h.observer.Logger.Error().Err(err).Msg("user is deleted")
	//	return domain.ErrUserIsDeleted
	//}

	if cmd.FirstName != nil {
		firstName, err := value_object.NewFirstName(*cmd.FirstName)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to add first name")
			return fmt.Errorf("failed to add first name: %w", err)
		}
		err = user.ChangeFirstName(firstName)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to change first name")
			return fmt.Errorf("failed to change first name: %w", err)
		}
	}

	if cmd.LastName != nil {
		lastName, err := value_object.NewLastName(*cmd.LastName)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to add last name")
			return fmt.Errorf("failed to add last name: %w", err)
		}
		err = user.ChangeLastName(lastName)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to change last name")
			return fmt.Errorf("failed to change last name: %w", err)
		}
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
