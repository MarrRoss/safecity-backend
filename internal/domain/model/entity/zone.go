package entity

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type Zone struct {
	ID         value_object.ID
	Name       value_object.ZoneName
	Boundaries *[]value_object.Location
	Safety     value_object.Safety
	CreatedAt  time.Time
	UpdatedAt  time.Time
	EndedAt    *time.Time
	FamilyID   value_object.ID
}

func NewZone(
	name value_object.ZoneName,
	bn *[]value_object.Location,
	safe value_object.Safety,
	familyID value_object.ID,
) (*Zone, error) {
	if bn == nil {
		return nil, fmt.Errorf("boundaries are nil: %w", domain.ErrInvalidZoneBoundaries)
	}
	id := value_object.NewID()
	newZone := Zone{
		ID:         id,
		Name:       name,
		Boundaries: bn,
		Safety:     safe,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		EndedAt:    nil,
		FamilyID:   familyID}
	return &newZone, nil
}

func (z *Zone) ChangeName(name value_object.ZoneName) error {
	if z.EndedAt != nil {
		return domain.ErrZoneIsDeleted
	}
	z.Name = name
	z.UpdatedAt = time.Now()
	return nil
}

func (z *Zone) ChangeSafety(safety value_object.Safety) error {
	if z.EndedAt != nil {
		return domain.ErrZoneIsDeleted
	}
	z.Safety = safety
	z.UpdatedAt = time.Now()
	return nil
}

func (z *Zone) ChangeBoundaries(boundaries *[]value_object.Location) error {
	if z.EndedAt != nil {
		return domain.ErrZoneIsDeleted
	}
	if boundaries == nil {
		return domain.ErrInvalidZoneBoundaries
	}
	z.Boundaries = boundaries
	z.UpdatedAt = time.Now()
	return nil
}

func (z *Zone) StopExistence() error {
	if z.EndedAt != nil {
		return domain.ErrZoneIsDeleted
	}
	timeNow := time.Now()
	z.EndedAt = &timeNow
	z.UpdatedAt = time.Now()
	return nil
}
