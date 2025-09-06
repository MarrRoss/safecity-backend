package user

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
)

type GetAllUsersQuery struct {
	Login *string
}

type GetAllUsersHandler struct {
	storage  db.UserRepository
	observer *observability.Observability
}

func NewGetAllUsersHandler(
	storage db.UserRepository,
	observer *observability.Observability,
) *GetAllUsersHandler {
	return &GetAllUsersHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *GetAllUsersHandler) Handle(ctx context.Context, qry GetAllUsersQuery) ([]*app_model.ApplicationUser, error) {
	var login *value_object.Login
	if qry.Login != nil {
		login1, err := value_object.NewLogin(*qry.Login)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid login format")
			return nil, fmt.Errorf("invalid login format: %w", err)
		}
		login = &login1
	}
	users, err := h.storage.GetUsers(ctx, login)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting users")
			return nil, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting users")
		return nil, err
	}

	return app_model.NewApplicationUsers(users), nil
}
