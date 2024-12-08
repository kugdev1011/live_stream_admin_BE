package model

import (
	"time"

	"gorm.io/gorm"
)

type AdminAction string

const (
	Login  AdminAction = "end_stream"
	Logout AdminAction = "banned_user"
)

type RoleType string

const (
	ADMINROLE       RoleType = "admin"
	USERROLE        RoleType = "user"
	GUESTROLE       RoleType = "guest"
	SUPPERADMINROLE RoleType = "supper_admin"
)

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Type        string    `gorm:"type:varchar(50);not null;unique" json:"type,omitempty"`
	Description string    `gorm:"type:text" json:"desription,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at,omitempty"`
	Users       []User    `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"users,omitempty"`
}

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Username     string         `gorm:"type:varchar(50);not null;unique" json:"user_name,omitempty"`
	DisplayName  string         `gorm:"type:varchar(100)" json:"display_name,omitempty"`
	Email        string         `gorm:"type:varchar(100);not null;unique" json:"email,omitempty"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	RoleID       uint           `gorm:"not null" json:"role_id,omitempty"`
	Role         *Role          `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	CreatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	CreatedByID  uint           `gorm:"index" json:"created_by_id,omitempty"`
	CreatedBy    *User          `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	UpdatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedByID  uint           `gorm:"index" json:"updated_by_id,omitempty"`
	UpdatedBy    *User          `gorm:"foreignKey:UpdatedByID" json:"updated_by,omitempty"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty"`
	DeletedByID  *uint          `json:"deleted_by_id,omitempty"`
	AdminLogs    []AdminLog     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"admin_logs,omitempty"`
}

type AdminLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	UserID      uint      `gorm:"not null" json:"user_id,omitempty"`
	Action      string    `gorm:"type:varchar(100);not null" json:"action,omitempty"`
	Details     string    `gorm:"type:text" json:"details,omitempty"`
	PerformedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"performed_at,omitempty"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type BlockedList struct {
	UserID        uint      `gorm:"primaryKey"`
	BlockedUserID uint      `gorm:"primaryKey"`
	User          User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	BlockedUser   User      `gorm:"foreignKey:BlockedUserID;constraint:OnDelete:CASCADE"`
	BlockedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
