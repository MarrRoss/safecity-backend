package notification_setting

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type GetAvailableNotificationSendersQuery struct {
	ReceiverID string
	FamilyID   uuid.UUID
}

type GetAvailableNotificationSendersHandler struct {
	userStorage       db.UserRepository
	familyStorage     db.FamilyRepository
	membershipStorage db.FamilyMembershipRepository
	observer          *observability.Observability
}

func NewGetAvailableNotificationSendersHandler(
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	membershipStorage db.FamilyMembershipRepository,
	observer *observability.Observability,
) *GetAvailableNotificationSendersHandler {
	return &GetAvailableNotificationSendersHandler{
		userStorage:       userStorage,
		familyStorage:     familyStorage,
		membershipStorage: membershipStorage,
		observer:          observer,
	}
}

func (h *GetAvailableNotificationSendersHandler) Handle(ctx context.Context,
	query GetAvailableNotificationSendersQuery) ([]*app_model.ApplicationUser, error) {
	receiverID, err := value_object.NewIDFromString(query.ReceiverID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return nil, fmt.Errorf("invalid user external id: %w", err)
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, receiverID)
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

	familyID, err := value_object.NewIDFromString(query.FamilyID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return nil, fmt.Errorf("invalid family id: %w", err)
	}
	family, err := h.familyStorage.GetFamily(ctx, familyID)
	if err != nil {
		if errors.Is(err, adapter.ErrFamilyNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("family not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
		return nil, err
	}

	senders, err := h.membershipStorage.GetAvailableNotificationSenders(ctx, user.ID, family.ID)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get available notification senders")
		return nil, fmt.Errorf("failed to get available notification senders: %w", err)
	}

	return app_model.NewApplicationUsers(senders), nil
}
