package family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
)

type AddFamilyCommand struct {
	Name   string
	UserID string
}

type AddFamilyHandler struct {
	familyStorage     db.FamilyRepository
	userStorage       db.UserRepository
	membershipStorage db.FamilyMembershipRepository
	roleStorage       db.RoleRepository
	observer          *observability.Observability
}

func NewAddFamilyHandler(familyStorage db.FamilyRepository,
	userStorage db.UserRepository,
	membershipStorage db.FamilyMembershipRepository,
	roleStorage db.RoleRepository,
	observer *observability.Observability) *AddFamilyHandler {
	return &AddFamilyHandler{
		familyStorage:     familyStorage,
		userStorage:       userStorage,
		membershipStorage: membershipStorage,
		roleStorage:       roleStorage,
		observer:          observer,
	}
}

func (h *AddFamilyHandler) Handle(ctx context.Context, cmd AddFamilyCommand) (value_object.ID, value_object.ID, error) {
	extID, err := value_object.NewIDFromString(cmd.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid external id")
		return value_object.ID{}, value_object.ID{}, err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("author not found")
			return value_object.ID{}, value_object.ID{}, err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting author")
			return value_object.ID{}, value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting author")
		return value_object.ID{}, value_object.ID{}, err
	}
	if user.EndedAt != nil {
		h.observer.Logger.Trace().Err(err).Msg("author is deleted")
		return value_object.ID{}, value_object.ID{}, domain.ErrUserIsDeleted
	}
	name, err := value_object.NewFamilyName(cmd.Name)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create family name")
		return value_object.ID{}, value_object.ID{}, fmt.Errorf("failed to create family name: %w", err)
	}
	family, err := entity.NewFamily(name)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create family")
		return value_object.ID{}, value_object.ID{}, fmt.Errorf("failed to create family: %w", err)
	}
	err = h.familyStorage.AddFamily(ctx, family)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding family")
			return value_object.ID{}, value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding family")
		return value_object.ID{}, value_object.ID{}, err
	}

	role, err := h.roleStorage.GetRole(ctx, 1)
	if err != nil {
		if errors.Is(err, adapter.ErrRoleNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("role not found")
			return value_object.ID{}, value_object.ID{}, err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting role")
			return value_object.ID{}, value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting role")
		return value_object.ID{}, value_object.ID{}, err
	}

	membership, err := aggregate.NewFamilyMembership(user, role, family)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create family membership")
		return value_object.ID{}, value_object.ID{}, fmt.Errorf("failed to create family membership: %w", err)
	}
	err = h.membershipStorage.AddFamilyMembership(ctx, membership)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding membership")
			return value_object.ID{}, value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding membership")
		return value_object.ID{}, value_object.ID{}, err
	}
	return family.ID, membership.ID, err
}
