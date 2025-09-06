package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type RoleDB struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func RoleDbToEntity(db *RoleDB) (*entity.Role, error) {
	roleName, err := value_object.NewRoleName(db.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to role name")
	}
	return &entity.Role{
		ID:        db.ID,
		Name:      roleName,
		CreatedAt: db.CreatedAt,
	}, nil
}

func RoleDBListToEntityList(dbList []*RoleDB) ([]*entity.Role, error) {
	var roles []*entity.Role
	for _, db := range dbList {
		role, err := RoleDbToEntity(db)
		if err != nil {
			return nil, fmt.Errorf("failed to convert storage data to entity")
		}
		roles = append(roles, role)
	}
	return roles, nil
}
