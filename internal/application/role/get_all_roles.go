package role

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetAllRolesHandler struct {
	roleStorage db.RoleRepository
}

func NewGetAllRolesHandler(roleStorage db.RoleRepository) *GetAllRolesHandler {
	return &GetAllRolesHandler{
		roleStorage: roleStorage,
	}
}

func (h *GetAllRolesHandler) Handle(ctx context.Context) ([]*app_model.ApplicationRole, error) {
	roles, err := h.roleStorage.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return app_model.NewApplicationRoles(roles), nil
}
