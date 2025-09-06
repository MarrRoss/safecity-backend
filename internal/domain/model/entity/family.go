package entity

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"time"
)

type Family struct {
	ID        value_object.ID
	Name      value_object.FamilyName
	CreatedAt time.Time
	UpdatedAt time.Time
	EndedAt   *time.Time
	Zones     *[]Zone
}

func NewFamily(name value_object.FamilyName) (*Family, error) {
	id := value_object.NewID()
	now := time.Now()
	newFamily := Family{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		EndedAt:   nil,
		Zones:     nil,
	}
	return &newFamily, nil
}

func (f *Family) ChangeName(name value_object.FamilyName) error {
	if f.EndedAt != nil {
		return domain.ErrFamilyIsDeleted
	}
	f.Name = name
	f.UpdatedAt = time.Now()
	return nil
}

func (f *Family) StopExistence() error {
	if f.EndedAt != nil {
		return domain.ErrFamilyIsDeleted
	}
	timeNow := time.Now()
	f.EndedAt = &timeNow
	f.UpdatedAt = time.Now()
	return nil
}
