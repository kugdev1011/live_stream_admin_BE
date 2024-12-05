package repository

import "gorm.io/gorm"

type Repository struct {
	User  *UserRepository
	Admin *AdminRepository
}

func NewRepository(db *gorm.DB) *Repository {
	adminRepo := newAdminRepository(db)
	userRepo := newUserRepository(db)
	return &Repository{
		Admin: adminRepo,
		User:  userRepo,
	}
}
