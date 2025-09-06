package app_model

import "awesomeProjectDDD/internal/domain/model/entity"

type ApplicationRole struct {
	ID        string
	Name      string
	CreatedAt string
}

func NewApplicationRole(role *entity.Role) *ApplicationRole {
	return &ApplicationRole{
		ID:        role.IDToString(),
		Name:      role.Name.String(),
		CreatedAt: role.CreatedAt.String(),
	}
}

func NewApplicationRoles(roles []*entity.Role) []*ApplicationRole {
	appRoles := make([]*ApplicationRole, len(roles))
	for i, role := range roles {
		appRoles[i] = NewApplicationRole(role)
	}
	return appRoles
}
