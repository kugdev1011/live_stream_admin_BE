package repository

import "gorm.io/gorm"

type Repository struct {
	User  *UserRepository
	Admin *AdminLogRepository
	Role  *RoleRepository
}

func NewRepository(db *gorm.DB) *Repository {
	adminRepo := newAdminRepository(db)
	userRepo := newUserRepository(db)
	roleRepo := NewRoleRepository(db)
	return &Repository{
		Admin: adminRepo,
		User:  userRepo,
		Role:  roleRepo,
	}
}
