package family

import (
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type DeleteFamilyQuery struct {
	ID string
}

type DeleteFamilyHandler struct {
	familiesStorage db.FamilyRepository
}

func NewDeleteFamilyHandler(familiesStorage db.FamilyRepository) *DeleteFamilyHandler {
	return &DeleteFamilyHandler{
		familiesStorage: familiesStorage,
	}
}

func (h *DeleteFamilyHandler) Handle(ctx context.Context, query DeleteFamilyQuery) error {
	id, err := value_object.NewIDFromString(query.ID)
	if err != nil {
		return fmt.Errorf("invalid family id: %w", err)
	}

	//memberships, err := h.membershipsStore.GetMembershipsByFamilyID(id)
	//if err != nil {
	//	return fmt.Errorf("failed to get memberships: %w", err)
	//}
	//for _, membership := range memberships {
	//	membershipID, err := value_object.NewIDFromString(membership.externalID)
	//	if err != nil {
	//		return fmt.Errorf("invalid membership id: %s", membership.externalID)
	//	}
	//	err = h.membershipsStore.DeleteMembership(membershipID)
	//	if err != nil {
	//		return fmt.Errorf("failed to delete membership: %w", err)
	//	}
	//}

	_, err = h.familiesStorage.GetFamily(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get family: %w", err)
	}
	//err = h.familyZoneStorage.DeleteFamilyZones(ctx, family.externalID)
	//if err != nil {
	//	return fmt.Errorf("failed to delete family zones: %w", err)
	//}

	// заменить на update
	//err = h.familiesStorage.DeleteFamily(ctx, family.externalID)
	//if err != nil {
	//	return fmt.Errorf("failed to delete family from storage: %w", err)
	//}
	return nil
}
