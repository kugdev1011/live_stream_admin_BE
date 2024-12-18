package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type AdminAction string

const (
	Login        AdminAction = "end_stream"
	Logout       AdminAction = "banned_user"
	LoginAction  AdminAction = "login"
	LogoutAction AdminAction = "logout"
	// Adjusted the constants to match their names
)

type RoleType string

const (
	SUPPERADMINROLE RoleType = "super_admin"
	ADMINROLE       RoleType = "admin"
	STREAMER        RoleType = "streamer"
	USERROLE        RoleType = "user"
)

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Type        RoleType  `gorm:"type:varchar(50);not null;unique" json:"type,omitempty"`
	Description string    `gorm:"type:text" json:"desription,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at,omitempty"`
	Users       []User    `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"users,omitempty"`
}

type User struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement"`
	Username            string         `gorm:"type:varchar(50);not null;unique"`
	DisplayName         string         `gorm:"type:varchar(100)" json:"display_name,omitempty"`
	Email               string         `gorm:"type:varchar(100);not null;unique"`
	PasswordHash        string         `gorm:"type:varchar(255);not null"`
	OTP                 string         `gorm:"type:varchar(6);null"`
	OTPExpiresAt        *time.Time     `gorm:"type:timestamp;null" json:"otp_expires_at,omitempty"`
	RoleID              uint           `gorm:"not null"`
	Role                Role           `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	CreatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	CreatedByID         *uint          `gorm:"index;null" json:"created_by_id,omitempty"`
	CreatedBy           *User          `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	UpdatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedByID         *uint          `gorm:"index;null" json:"updated_by_id,omitempty"`
	UpdatedBy           *User          `gorm:"foreignKey:UpdatedByID" json:"updated_by,omitempty"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty"`
	DeletedByID         *uint          `json:"deleted_by_id,omitempty"`
	AvatarFileName      sql.NullString `gorm:"type:varchar(255)" json:"avatar_file_name,omitempty"`
	AdminLogs           []AdminLog     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedByCategories []Category     `gorm:"foreignKey:CreatedByID"`
	UpdatedByCategories []Category     `gorm:"foreignKey:UpdatedByID"`
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

type TwoFA struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"not null;unique"`
	Secret       string    `gorm:"type:text;not null"`
	Is2faEnabled bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
