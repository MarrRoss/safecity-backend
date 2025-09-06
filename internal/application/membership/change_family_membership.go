package membership

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"github.com/google/uuid"
)

type ChangeFamilyMembershipCommand struct {
	ID     uuid.UUID
	RoleID int
}

type ChangeFamilyMembershipHandler struct {
	familyStorage     db.FamilyRepository
	membershipStorage db.FamilyMembershipRepository
	observer          *observability.Observability
}

func NewChangeFamilyMembershipHandler(
	familyStorage db.FamilyRepository,
	membershipStorage db.FamilyMembershipRepository,
	observer *observability.Observability) *ChangeFamilyMembershipHandler {
	return &ChangeFamilyMembershipHandler{
		familyStorage:     familyStorage,
		membershipStorage: membershipStorage,
		observer:          observer,
	}
}

func (h *ChangeFamilyMembershipHandler) Handle(ctx context.Context, cmd ChangeFamilyMembershipCommand) error {
	membershipID, err := value_object.NewIDFromString(cmd.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid membership id")
		return err
	}
	err = h.membershipStorage.UpdateFamilyMembershipRole(ctx, membershipID, cmd.RoleID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while updating role in membership")
			return err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while updating role in membership")
		return err
	}

	//family, err := h.storage.GetFamily(ctx, id)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrFamilyNotFound) {
	//		h.observer.Logger.Trace().Err(err).Msg("family not found")
	//		return err
	//	}
	//
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting family")
	//		return err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
	//	return err
	//}
	//
	//name, err := value_object.NewFamilyName(cmd.Name)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to create family name")
	//	return fmt.Errorf("failed to create family name: %w", err)
	//}
	//err = family.ChangeName(name)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to change family name")
	//	return fmt.Errorf("failed to change family name: %w", err)
	//}
	//err = h.storage.UpdateFamily(ctx, family)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while updating family")
	//		return err
	//	}
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while updating family")
	//	return err
	//}
	return nil
}
