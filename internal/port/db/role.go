package db

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"context"
)

type RoleRepository interface {
	AddRole(ctx context.Context, role *entity.Role) error
	GetRole(ctx context.Context, id int) (*entity.Role, error)
	GetRoles(ctx context.Context) ([]*entity.Role, error)
}
