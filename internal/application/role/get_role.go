package role

import (
	"awesomeProjectDDD/internal/adapter/hydra"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
	"strconv"
)

type GetRoleQuery struct {
	ID             string
	ExternalUserID string
}

type GetRoleHandler struct {
	roleStorage db.RoleRepository
	ssoService  *hydra.Service
}

func NewGetRoleHandler(roleStorage db.RoleRepository, ssoService *hydra.Service) *GetRoleHandler {
	return &GetRoleHandler{
		roleStorage: roleStorage,
		ssoService:  ssoService,
	}
}

func (h *GetRoleHandler) Handle(ctx context.Context, query GetRoleQuery) (*app_model.ApplicationRole, error) {
	roleID, err := strconv.Atoi(query.ID)
	if err != nil {
		return &app_model.ApplicationRole{}, fmt.Errorf("invalid role id: %w", err)
	}
	role, err := h.roleStorage.GetRole(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return app_model.NewApplicationRole(role), nil
}
