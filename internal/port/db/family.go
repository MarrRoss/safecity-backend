package db

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type FamilyRepository interface {
	AddFamily(ctx context.Context, family *entity.Family) error
	GetFamily(ctx context.Context, id value_object.ID) (*entity.Family, error)
	GetFamilyZones(ctx context.Context, id value_object.ID) ([]*entity.Zone, error)
	//GetFamiliesByAuthorID(ctx context.Context, id value_object.ID) ([]*entity.Family, error)
	GetFamiliesByIDs(ctx context.Context, ids []string) ([]*entity.Family, error)
	UpdateFamily(ctx context.Context, family *entity.Family) error
	//DeleteFamily(ctx context.Context, id value_object.externalID) error
}
