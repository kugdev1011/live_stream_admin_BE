package dto

import (
	"gitlab/live/be-live-api/model"
	"time"
)

type CreateAdminRequest struct {
	UserName    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email,max=100"`
	DisplayName string `json:"display_name" validate:"required,min=3,max=100"`
	Password    string `json:"password" validate:"required,min=6,max=255"`
	CreatedByID *uint  `json:"created_by_id" validate:"required"`
}

type CreateAdminResp struct {
	UserName    string         `json:"username,omitempty"`
	DisplayName string         `json:"display_name,omitempty"`
	Email       string         `json:"email,omitempty"`
	Role        model.RoleType `json:"role,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
}
