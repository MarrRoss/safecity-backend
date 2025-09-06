package user

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AddUserCommand struct {
	FirstName  string
	LastName   string
	Email      string
	Login      string
	ExternalID uuid.UUID
	Tracking   bool
}

type AddUserHandler struct {
	storage  db.UserRepository
	observer *observability.Observability
}

func NewAddUserHandler(storage db.UserRepository, observer *observability.Observability) *AddUserHandler {
	return &AddUserHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *AddUserHandler) Handle(ctx context.Context, cmd AddUserCommand) (value_object.ID, error) {
	email, err := value_object.NewEmail(cmd.Email)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add email")
		return value_object.ID{}, fmt.Errorf("failed to add email: %w", err)
	}
	exists, err := h.storage.EmailExists(ctx, email)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting email exists")
			return value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting email exists")
		return value_object.ID{}, err
	}
	if exists {
		h.observer.Logger.Trace().Msg("this email already exists in system")
		return value_object.ID{}, application.ErrEmailExists
	}

	firstName, err := value_object.NewFirstName(cmd.FirstName)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add first name")
		return value_object.ID{}, fmt.Errorf("failed to add first name: %w", err)
	}
	lastName, err := value_object.NewLastName(cmd.LastName)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add last name")
		return value_object.ID{}, fmt.Errorf("failed to add last name: %w", err)
	}
	name := value_object.NewFullName(firstName, lastName)
	login, err := value_object.NewLogin(cmd.Login)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add login")
		return value_object.ID{}, fmt.Errorf("failed to add login: %w", err)
	}
	extID, err := value_object.NewIDFromString(cmd.ExternalID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add external id")
		return value_object.ID{}, fmt.Errorf("failed to add external id: %w", err)
	}
	user, err := entity.NewUser(name, email, login, extID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create user")
		return value_object.ID{}, fmt.Errorf("failed to create user: %w", err)
	}
	err = h.storage.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding user")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding user")
		return value_object.ID{}, err
	}
	return user.ID, nil
}
