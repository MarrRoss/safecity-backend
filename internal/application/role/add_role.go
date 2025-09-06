package role

import (
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type AddRoleCommand struct {
	Name string
}

type AddRoleHandler struct {
	roleStorage db.RoleRepository
}

func NewAddRoleHandler(roleStorage db.RoleRepository) *AddRoleHandler {
	return &AddRoleHandler{
		roleStorage: roleStorage,
	}
}

func (h *AddRoleHandler) Handle(ctx context.Context, cmd AddRoleCommand) (string, error) {
	//name, err := value_object.NewRoleName(cmd.Name)
	//if err != nil {
	//	return "", fmt.Errorf("failed to create role name: %w", err)
	//}
	//role, err := entity.NewRole(name)
	//if err != nil {
	//	return "", fmt.Errorf("failed to create role: %w", err)
	//}
	//err = h.roleStorage.AddRole(ctx, role)
	//if err != nil {
	//	return "", fmt.Errorf("failed to add role to storage: %w", err)
	//}
	//return role.IDToString(), err
	return "", nil
}
