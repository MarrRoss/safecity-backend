package membership

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

type GetUserMembershipsQuery struct {
	ID string
}

type GetUserMembershipsHandler struct {
	membershipStorage db.FamilyMembershipRepository
	userStorage       db.UserRepository
	familyStorage     db.FamilyRepository
	observer          *observability.Observability
}

func NewGetUserMembershipsHandler(
	membershipStorage db.FamilyMembershipRepository,
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	observer *observability.Observability) *GetUserMembershipsHandler {
	return &GetUserMembershipsHandler{
		membershipStorage: membershipStorage,
		userStorage:       userStorage,
		familyStorage:     familyStorage,
		observer:          observer,
	}
}

func (h *GetUserMembershipsHandler) Handle(
	ctx context.Context,
	query GetUserMembershipsQuery) (
	[]*app_model.ApplicationFamilyMembership, error) {
	extID, err := value_object.NewIDFromString(query.ID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid external id")
		return nil, err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
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

	memberships, err := h.membershipStorage.GetMembershipsParticipantsByUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrMembershipsNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user memberships not found")
			return nil, err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user memberships")
			return nil, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user memberships")
		return nil, err
	}

	var userIDs []string
	var familyIDs []string
	for _, membership := range memberships {
		userIDs = append(userIDs, membership.UserID)
		familyIDs = append(familyIDs, membership.FamilyID)
	}
	users, err := h.userStorage.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrMembershipUsersNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("membership users not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting membership users")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting membership users")
		return nil, err
	}
	families, err := h.familyStorage.GetFamiliesByIDs(ctx, familyIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrMembershipFamiliesNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("membership families not found")
			return nil, err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting membership families")
			return nil, err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting membership families")
		return nil, err
	}

	resMemberships, err := app_model.NewApplicationFamilyMemberships(user.ID.String(), memberships, users, families)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to create new application family memberships")
		return nil, fmt.Errorf("failed to create new application family memberships: %w", err)
	}
	//resMemberships, err := app_model.NewApplicationFamilyMemberships(memberships, users, families)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to create new application family memberships")
	//	return nil, fmt.Errorf("failed to create new application family memberships: %w", err)
	//}
	return resMemberships, nil
}
