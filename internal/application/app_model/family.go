package app_model

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"github.com/google/uuid"
)

type ApplicationFamily struct {
	ID   uuid.UUID
	Name string
}

func NewApplicationFamily(family *entity.Family) *ApplicationFamily {
	return &ApplicationFamily{
		ID:   family.ID.ToRaw(),
		Name: family.Name.String(),
	}
}

func NewApplicationFamilies(families []*entity.Family) []*ApplicationFamily {
	appFamilies := make([]*ApplicationFamily, len(families))
	for i, family := range families {
		appFamilies[i] = NewApplicationFamily(family)
	}
	return appFamilies
}
