package repository

import (
	"gitlab/live/be-live-api/model"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func (s *AdminRepository) CreateAdmin(newUser *model.User) (*model.User, error) {

	// find role admin
	var role model.Role
	if err := s.db.Model(model.Role{}).Where("type=?", model.ADMINROLE).First(&role).Error; err != nil {
		return nil, err
	}

	newUser.RoleID = role.ID
	if err := s.db.Model(model.User{}).Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *AdminRepository) ById(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.Model(model.User{}).Where("id=? AND deleted_at IS NULL", id).
		Preload("Role").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("AdminLogs").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AdminRepository) Create(adminLog *model.AdminLog) error {
	if err := r.db.Create(adminLog).Error; err != nil {
		return err
	}
	return nil
}
func newAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}
