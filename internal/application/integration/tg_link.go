package integration

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
)

type AddTgLinkCommand struct {
	UserID string
}

type AddTgLinkHandler struct {
	userStorage        db.UserRepository
	integrationStorage db.IntegrationRepository
	observer           *observability.Observability
}

func NewAddTgLinkHandler(userStorage db.UserRepository,
	integrationStorage db.IntegrationRepository,
	observer *observability.Observability) *AddTgLinkHandler {
	return &AddTgLinkHandler{
		userStorage:        userStorage,
		integrationStorage: integrationStorage,
		observer:           observer,
	}
}

func (h *AddTgLinkHandler) Handle(ctx context.Context, cmd AddTgLinkCommand) (string, error) {
	extID, err := value_object.NewIDFromString(cmd.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid external id")
		return "", err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return "", err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return "", err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return "", err
	}
	integration, err := h.integrationStorage.GetIntegration(ctx, user.ID, 1)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting integration")
			return "", err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting integration")
		return "", err
	}

	if !user.TgIntegrated {
		if integration != nil {
			fmt.Printf("integration: %v\n", integration.ID)
			link := fmt.Sprintf("https://t.me/safecityrobot?start=%s", user.ID)
			h.observer.Logger.Info().Msg(link)
			return link, nil
		} else {
			fmt.Printf("integration is nil")
			err := h.integrationStorage.AddIntegration(ctx, user.ID, 1)
			if err != nil {
				if errors.Is(err, adapter.ErrStorage) {
					h.observer.Logger.Error().Err(err).Msg("database error while adding integration")
					return "", adapter.ErrIntegrationAlreadyExists
				}
				h.observer.Logger.Error().Err(err).Msg("unexpected error while adding integration")
				return "", err
			}
			link := fmt.Sprintf("https://t.me/safecityrobot?start=%s", user.ID)
			h.observer.Logger.Info().Msg(link)
			return link, nil
		}
	}
	h.observer.Logger.Trace().Err(err).Msg("user is already integrated")
	return "", domain.ErrUserAlreadyIntegrated
}
