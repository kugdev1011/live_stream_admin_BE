package dto

import "gitlab/live/be-live-api/model"

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type RegisterDTO struct {
	Username string         `json:"username" validate:"required,min=3,max=50"`
	Email    string         `json:"email" validate:"required,email,max=100"`
	Password string         `json:"password" validate:"required,min=6,max=255"`
	RoleType model.RoleType `json:"roleType" validate:"required,oneof=admin user guest"`
}

type ForgetPasswordDTO struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

type ResetPasswordDTO struct {
	OTP             string `json:"otp" validate:"required,len=6"`
	NewPassword     string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}
