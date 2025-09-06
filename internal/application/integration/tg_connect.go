package integration

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"github.com/google/uuid"
)

type AddTgConnectCommand struct {
	UserID     uuid.UUID
	TelegramID string
}

type AddTgConnectHandler struct {
	userStorage        db.UserRepository
	integrationStorage db.IntegrationRepository
	observer           *observability.Observability
}

func NewAddTgConnectHandler(userStorage db.UserRepository,
	integrationStorage db.IntegrationRepository,
	observer *observability.Observability) *AddTgConnectHandler {
	return &AddTgConnectHandler{
		userStorage:        userStorage,
		integrationStorage: integrationStorage,
		observer:           observer,
	}
}

func (h *AddTgConnectHandler) Handle(ctx context.Context, cmd AddTgConnectCommand) error {
	id, err := value_object.NewIDFromString(cmd.UserID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user id")
		return err
	}
	user, err := h.userStorage.GetUser(ctx, id)
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
	if user.EndedAt != nil {
		h.observer.Logger.Trace().Err(err).Msg("user is deleted")
		return domain.ErrUserIsDeleted
	}
	if !user.TgIntegrated {
		err := h.integrationStorage.UpdateIntegration(ctx, user.ID, cmd.TelegramID)
		if err != nil {
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while adding family")
				return err
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while adding family")
			return err
		}
		h.observer.Logger.Info().Msg("telegram connected")
		return nil
	}
	h.observer.Logger.Trace().Err(err).Msg("failed to connect telegram")
	return domain.ErrTgConnect
}
