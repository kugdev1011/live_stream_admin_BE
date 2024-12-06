package repository

import (
	"gitlab/live/be-live-api/model"
	"gorm.io/gorm"
)

type AdminLogRepository struct {
	db *gorm.DB
}

func newAdminRepository(db *gorm.DB) *AdminLogRepository {
	return &AdminLogRepository{
		db: db,
	}
}

func (r *AdminLogRepository) Create(adminLog *model.AdminLog) error {
	if err := r.db.Create(adminLog).Error; err != nil {
		return err
	}
	return nil
}
