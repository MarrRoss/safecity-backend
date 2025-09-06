package request

type GetRoleByIDRequest struct {
	ID string `json:"id" path:"id"`
}

type AddRoleRequest struct {
	Name string `json:"name"`
}
