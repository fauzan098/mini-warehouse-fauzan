package response

type LoginResponse struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"Email"`
	Role   []string `json:"role_name"`
}
