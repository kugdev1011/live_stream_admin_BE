package repository

import (
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/pkg/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (s *UserRepository) Page(page, limit int) (*utils.PaginationModel[model.User], error) {
	var query = s.db.Model(model.User{}).Preload("Role").Preload("AdminLogs")
	return utils.CreatePage[model.User](query, page, limit)
}

func newUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
