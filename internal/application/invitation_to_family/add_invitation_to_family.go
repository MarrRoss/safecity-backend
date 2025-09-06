package invitation_to_family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AddInvitationToFamilyCommand struct {
	AuthorID string
	RoleID   int
	FamilyID uuid.UUID
}

type AddInvitationToFamilyHandler struct {
	invitationStorage db.InvitationToFamilyRepository
	userStorage       db.UserRepository
	roleStorage       db.RoleRepository
	familyStorage     db.FamilyRepository
	observer          *observability.Observability
}

func NewAddInvitationToFamilyHandler(
	invitationStorage db.InvitationToFamilyRepository,
	userStorage db.UserRepository,
	roleStorage db.RoleRepository,
	familyStorage db.FamilyRepository,
	observer *observability.Observability,
) *AddInvitationToFamilyHandler {
	return &AddInvitationToFamilyHandler{
		invitationStorage: invitationStorage,
		userStorage:       userStorage,
		roleStorage:       roleStorage,
		familyStorage:     familyStorage,
		observer:          observer,
	}
}

func (h *AddInvitationToFamilyHandler) Handle(ctx context.Context,
	cmd AddInvitationToFamilyCommand) (value_object.ID, error) {
	extID, err := value_object.NewIDFromString(cmd.AuthorID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid external id")
		return value_object.ID{}, err
	}
	familyID, err := value_object.NewIDFromString(cmd.FamilyID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return value_object.ID{}, err
	}

	author, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("author not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting author")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting author")
		return value_object.ID{}, err
	}

	role, err := h.roleStorage.GetRole(ctx, cmd.RoleID)
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
	invitation, err := aggregate.NewInvitationToFamily(author, role, family)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create invitation to family")
		return value_object.ID{}, fmt.Errorf("failed to create invitation to family: %w", err)
	}
	err = h.invitationStorage.AddInvitationToFamily(ctx, invitation)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding invitation to family")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding invitation to family")
		return value_object.ID{}, err
	}
	return invitation.ID, nil
}
