package response

import "awesomeProjectDDD/internal/application/app_model"

type GetRoleByIDResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGetRoleByIDResponse(role *app_model.ApplicationRole) *GetRoleByIDResponse {
	return &GetRoleByIDResponse{
		ID:   role.ID,
		Name: role.Name,
	}
}

func NewGetRolesResponse(roles []*app_model.ApplicationRole) []*GetRoleByIDResponse {
	respRoles := make([]*GetRoleByIDResponse, len(roles))
	for i, role := range roles {
		respRoles[i] = NewGetRoleByIDResponse(role)
	}
	return respRoles
}
