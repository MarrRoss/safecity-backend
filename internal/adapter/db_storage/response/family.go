package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type FamilyDB struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func FamilyDbToEntity(db *FamilyDB) (*entity.Family, error) {
	familyID, err := value_object.NewIDFromString(db.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to family id")
	}
	familyName, err := value_object.NewFamilyName(db.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to family name: %v", err)
	}
	return &entity.Family{
		ID:        familyID,
		Name:      familyName,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}, nil
}

func FamilyDbListToEntityList(dbList []*FamilyDB) ([]*entity.Family, error) {
	var families []*entity.Family
	for _, db := range dbList {
		family, err := FamilyDbToEntity(db)
		if err != nil {
			return nil, fmt.Errorf("failed to convert storage data to entity: %v", err)
		}
		families = append(families, family)
	}
	return families, nil
}
