package membership

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AddUserMembershipCommand struct {
	ID     uuid.UUID
	UserID string
}

type AddUserMembershipHandler struct {
	membershipStorage db.FamilyMembershipRepository
	userStorage       db.UserRepository
	roleStorage       db.RoleRepository
	familyStorage     db.FamilyRepository
	invitationStorage db.InvitationToFamilyRepository
	observer          *observability.Observability
}

func NewAddUserMembershipHandler(
	membershipStorage db.FamilyMembershipRepository,
	userStorage db.UserRepository,
	roleStorage db.RoleRepository,
	familyStorage db.FamilyRepository,
	invitationStorage db.InvitationToFamilyRepository,
	observer *observability.Observability) *AddUserMembershipHandler {
	return &AddUserMembershipHandler{
		membershipStorage: membershipStorage,
		userStorage:       userStorage,
		roleStorage:       roleStorage,
		familyStorage:     familyStorage,
		invitationStorage: invitationStorage,
		observer:          observer,
	}
}

func (h *AddUserMembershipHandler) Handle(ctx context.Context, cmd AddUserMembershipCommand) (value_object.ID, error) {
	invitationID, err := value_object.NewIDFromString(cmd.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid invitation id")
		return value_object.ID{}, err
	}
	invitation, err := h.invitationStorage.GetInvitationToFamily(ctx, invitationID)
	if err != nil {
		if errors.Is(err, adapter.ErrInvitationNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("invitation not found")
			return value_object.ID{}, err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting invitation")
			return value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting invitation")
		return value_object.ID{}, err
	}
	if invitation.AuthorID == cmd.UserID {
		h.observer.Logger.Error().Err(err).Msg("user is invitation creator")
		return value_object.ID{}, domain.ErrInvalidID
	}
	extID, err := value_object.NewIDFromString(cmd.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user id")
		return value_object.ID{}, err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return value_object.ID{}, err
	}

	familyID, err := value_object.NewIDFromString(invitation.FamilyID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return value_object.ID{}, err
	}
	family, err := h.familyStorage.GetFamily(ctx, familyID)
	if err != nil {
		if errors.Is(err, adapter.ErrFamilyNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("family not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
		return value_object.ID{}, err
	}

	exists, err := h.membershipStorage.MembershipExists(ctx, user.ID, family.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while checking membership exists")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while checking membership exists")
		return value_object.ID{}, err
	}
	if exists {
		h.observer.Logger.Trace().Msg("this membership already exists in system")
		return value_object.ID{}, application.ErrMembershipExists
	}

	role, err := h.roleStorage.GetRole(ctx, invitation.RoleID)
	if err != nil {
		if errors.Is(err, adapter.ErrRoleNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("role not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting role")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting role")
		return value_object.ID{}, err
	}

	membership, err := aggregate.NewFamilyMembership(user, role, family)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create family membership")
		return value_object.ID{}, fmt.Errorf("failed to create family membership: %w", err)
	}
	err = h.membershipStorage.AddFamilyMembership(ctx, membership)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding membership")
			return value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding membership")
		return value_object.ID{}, err
	}

	err = h.invitationStorage.AddInviteActivation(ctx, invitationID, user.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding invite activation")
			return value_object.ID{}, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding invite activation")
		return value_object.ID{}, err
	}
	return membership.ID, nil
}
