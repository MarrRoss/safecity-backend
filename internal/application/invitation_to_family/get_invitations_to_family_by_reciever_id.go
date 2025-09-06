package invitation_to_family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"github.com/google/uuid"
)

type GetInvitationsToFamilyByReceiverIdQuery struct {
	//UserID       string
	InvitationID uuid.UUID
}

type GetInvitationsToFamilyByReceiverIdHandler struct {
	invitationStorage db.InvitationToFamilyRepository
	userStorage       db.UserRepository
	familyStorage     db.FamilyRepository
	roleStorage       db.RoleRepository
	observer          *observability.Observability
}

func NewGetInvitationsToFamilyByReceiverIdHandler(
	invitationStorage db.InvitationToFamilyRepository,
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	roleStorage db.RoleRepository,
	observer *observability.Observability,
) *GetInvitationsToFamilyByReceiverIdHandler {
	return &GetInvitationsToFamilyByReceiverIdHandler{
		invitationStorage: invitationStorage,
		userStorage:       userStorage,
		familyStorage:     familyStorage,
		roleStorage:       roleStorage,
		observer:          observer,
	}
}

func (h *GetInvitationsToFamilyByReceiverIdHandler) Handle(ctx context.Context,
	query GetInvitationsToFamilyByReceiverIdQuery) (*app_model.ApplicationInvitationToFamilyType, error) {
	//extID, err := value_object.NewIDFromString(query.UserID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
	//	return nil, err
	//}
	//user, err := h.userStorage.GetUserByExternalID(ctx, extID)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrUserNotFound) {
	//		h.observer.Logger.Trace().Err(err).Msg("user not found")
	//		return nil, err
	//	}
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting user")
	//		return nil, err
	//	}
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
	//	return nil, err
	//}
	invitationID, err := value_object.NewIDFromString(query.InvitationID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid invitation id")
		return nil, err
	}
	invitation, err := h.invitationStorage.GetInvitationToFamilyByID(ctx, invitationID)
	if err != nil {
		if errors.Is(err, adapter.ErrInvitationNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("invitation not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting invitation")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting invitation")
		return nil, err
	}
	//invitations, err := h.invitationStorage.GetPendingInvitationsByUser(ctx, user.ID)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting invitations")
	//		return nil, err
	//	}
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting invitations")
	//	return nil, err
	//}

	return app_model.NewApplicationInvitationToFamilyType(invitation), nil
}
