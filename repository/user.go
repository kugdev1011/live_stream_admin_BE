package repository

import (
	"gitlab/live/be-live-api/model"
	apimodel "gitlab/live/be-live-api/model/api-model"
	"gitlab/live/be-live-api/pkg/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (s *UserRepository) Page(filter *apimodel.UserQuery, page, limit int) (*utils.PaginationModel[model.User], error) {
	var query = s.db.Model(model.User{})
	if filter != nil && filter.Role != "" {
		query = query.Joins("LEFT JOIN roles ON roles.id = users.role_id").
			Where("roles.type = ?", filter.Role)
	}

	query = query.Preload("Role").Preload("AdminLogs").Preload("CreatedBy").Preload("UpdatedBy")
	pagination, err := utils.CreatePage[model.User](query, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.Create(pagination, page, limit)
}

func newUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
