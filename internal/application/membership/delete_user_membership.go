package membership

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

type DeleteUserMembershipCommand struct {
	ID     uuid.UUID
	UserID string
}

type DeleteUserMembershipHandler struct {
	membershipStorage db.FamilyMembershipRepository
	userStorage       db.UserRepository
	observer          *observability.Observability
}

func NewDeleteUserMembershipHandler(
	membershipStorage db.FamilyMembershipRepository,
	userStorage db.UserRepository,
	observer *observability.Observability,
) *DeleteUserMembershipHandler {
	return &DeleteUserMembershipHandler{
		membershipStorage: membershipStorage,
		userStorage:       userStorage,
		observer:          observer,
	}
}

func (h *DeleteUserMembershipHandler) Handle(ctx context.Context, cmd DeleteUserMembershipCommand) error {
	membershipID, err := value_object.NewIDFromString(cmd.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid membership id")
		return err
	}
	membership, err := h.membershipStorage.GetFamilyMembershipByID(ctx, membershipID)
	if err != nil {
		if errors.Is(err, adapter.ErrMembershipNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("membership not found")
			return err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting membership")
			return err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting membership")
		return err
	}

	extID, err := value_object.NewIDFromString(cmd.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid external id")
		return err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
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
	// user удаляет себя из семьи
	if membership.UserID == user.ID.String() {
		// TODO: проверить, остались ли в семье взрослые
		err = h.membershipStorage.DeleteFamilyMembership(ctx, membershipID)
		if err != nil {
			if errors.Is(err, adapter.ErrMembershipNotFound) {
				h.observer.Logger.Trace().Err(err).Msg("membership not found")
				return err
			}
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while updating membership")
				return err
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while updating membership")
			return err
		}
		fmt.Printf("Admin %v удален из семьи", membership.UserID)
		// user удаляет кого-то из семьи
	} else {
		err = h.membershipStorage.DeleteFamilyMembership(ctx, membershipID)
		if err != nil {
			if errors.Is(err, adapter.ErrMembershipNotFound) {
				h.observer.Logger.Trace().Err(err).Msg("membership not found")
				return err
			}
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while updating membership")
				return err
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while updating membership")
			return err
		}
		fmt.Printf("User %v удален из семьи", membership.UserID)
	}
	return nil
}
