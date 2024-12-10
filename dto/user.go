package dto

import (
	"time"

	"gorm.io/gorm"
)

type UserQuery struct {
	Role string `query:"role" validate:"omitempty,oneof=super_admin admin streamer user"`
}

type UserResponseDTO struct {
	ID          uint             `json:"id,omitempty"`
	Username    string           `json:"username,omitempty"`
	DisplayName string           `json:"display_name,omitempty"`
	Email       string           `json:"email,omitempty"`
	RoleID      uint             `json:"role_id,omitempty"`
	Role        *RoleDTO         `json:"role,omitempty"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
	CreatedByID *uint            `json:"created_by_id,omitempty"`
	CreatedBy   *UserResponseDTO `json:"created_by,omitempty"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
	UpdatedByID *uint            `json:"updated_by_id,omitempty"`
	UpdatedBy   *UserResponseDTO `json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at,omitempty"`
	DeletedByID *uint            `json:"deleted_by_id,omitempty"`
	AdminLogs   []AdminLogDTO    `json:"admin_logs"`
}

type RoleDTO struct {
	ID          uint              `json:"id,omitempty"`
	Type        string            `json:"type,omitempty"`
	Description string            `json:"description,omitempty"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
	Users       []UserResponseDTO `json:"users"`
}

type AdminLogDTO struct {
	ID          uint            `json:"id,omitempty"`
	UserID      uint            `json:"user_id,omitempty"`
	Action      string          `json:"action,omitempty"`
	Details     string          `json:"details,omitempty"`
	PerformedAt time.Time       `json:"performed_at,omitempty"`
	User        UserResponseDTO `json:"user,omitempty"`
}
