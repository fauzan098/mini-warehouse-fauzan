package request

type CreatRoleRequest struct {
	Name string `json:"name" validate:"required"`
}