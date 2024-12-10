package dto

type CreateAdminRequest struct {
	UserName    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email,max=100"`
	DisplayName string `json:"display_name" validate:"required,min=5,max=100"`
	Password    string `json:"password" validate:"required,min=6,max=255"`
	CreatedByID *uint  `json:"created_by_id" validate:"required"`
}
