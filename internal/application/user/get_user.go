package user

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
)

type GetUserQuery struct {
	ID string
}

type GetUserHandler struct {
	observer  *observability.Observability
	userStore db.UserRepository
}

func NewGetUserHandler(
	observer *observability.Observability,
	userStore db.UserRepository,
) *GetUserHandler {
	return &GetUserHandler{
		observer:  observer,
		userStore: userStore,
	}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*app_model.ApplicationUser, error) {
	extID, err := value_object.NewIDFromString(query.ID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return nil, err
	}
	user, err := h.userStore.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return nil, err
	}
	return app_model.NewApplicationUser(user), nil
}
