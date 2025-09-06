package membership

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
)

type GetMembershipsByFamilyQuery struct {
	FamilyID string
}

type GetMembershipsByFamilyHandler struct {
	membershipStorage db.FamilyMembershipRepository
	userStorage       db.UserRepository
	roleStorage       db.RoleRepository
	familyStorage     db.FamilyRepository
}

func NewGetMembershipsByFamilyHandler(
	membershipStorage db.FamilyMembershipRepository,
	userStorage db.UserRepository,
	roleStorage db.RoleRepository,
	familyStorage db.FamilyRepository) *GetMembershipsByFamilyHandler {
	return &GetMembershipsByFamilyHandler{
		membershipStorage: membershipStorage,
		userStorage:       userStorage,
		roleStorage:       roleStorage,
		familyStorage:     familyStorage,
	}
}

func (h *GetMembershipsByFamilyHandler) Handle(ctx context.Context,
	query GetMembershipsByFamilyQuery) ([]*response.FamilyMembershipDB, error) {
	familyID, err := value_object.NewIDFromString(query.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("invalid family id: %w", err)
	}
	memberships, err := h.membershipStorage.GetMembershipsByFamilyID(ctx, familyID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get family memberships: %w", err)
	}
	if len(memberships) == 0 {
		return nil, errors.New("memberships weren't found for family")
	}

	return memberships, nil
}
