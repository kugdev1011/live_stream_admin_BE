package model

import "time"

type AdminAction string

const (
	Login  AdminAction = "end_stream"
	Logout AdminAction = "banned_user"
)

const (
	LoginAction  AdminAction = "login"
	LogoutAction AdminAction = "logout"
	// Adjusted the constants to match their names
)

type RoleType string

const (
	ADMINROLE RoleType = "admin"
	USERROLE  RoleType = "user"
	GUESTROLE RoleType = "guest"
)

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Type        string    `gorm:"type:varchar(50);not null;unique"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Users       []User    `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	Username     string     `gorm:"type:varchar(50);not null;unique"`
	Email        string     `gorm:"type:varchar(100);not null;unique"`
	PasswordHash string     `gorm:"type:varchar(255);not null"`
	OTP          string     `gorm:"type:varchar(6);null"`
	OTPExpiresAt time.Time  `gorm:"type:timestamp;null"`
	RoleID       uint       `gorm:"not null"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
	AdminLogs    []AdminLog `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type AdminLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UserID      uint      `gorm:"not null"`
	Action      string    `gorm:"type:varchar(100);not null"`
	Details     string    `gorm:"type:text"`
	PerformedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
