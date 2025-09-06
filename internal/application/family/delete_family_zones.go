package family

import (
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type DeleteFamilyZonesCommand struct {
	FamilyID string
}

type DeleteFamilyZonesHandler struct {
	familyStorage db.FamilyRepository
}

func NewDeleteFamilyZonesHandler(
	familyStorage db.FamilyRepository) *DeleteFamilyZonesHandler {
	return &DeleteFamilyZonesHandler{
		familyStorage: familyStorage,
	}
}

func (h *DeleteFamilyZonesHandler) Handle(ctx context.Context, cmd DeleteFamilyZonesCommand) error {
	//familyID, err := value_object.NewIDFromString(cmd.FamilyID)
	//if err != nil {
	//	return fmt.Errorf("invalid family id: %w", err)
	//}
	//_, err = h.familyStorage.GetFamily(ctx, familyID)
	//if err != nil {
	//	return fmt.Errorf("failed to get family: %w", err)
	//}

	// заменить на update
	//err = h.familyZoneStorage.DeleteFamilyZones(ctx, family.externalID)
	//if err != nil {
	//	return fmt.Errorf("failed to delete family_zones from storage: %w", err)
	//}
	return nil
}
