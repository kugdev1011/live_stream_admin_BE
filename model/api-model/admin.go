package model

type CreateAdminRequest struct {
	UserName    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	CreatedByID uint   `json:"created_by_id"`
}
