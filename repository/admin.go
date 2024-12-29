package repository

import (
	"fmt"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/utils"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func (s *AdminRepository) CreateAdmin(newUser *model.User) error {

	// find role admin
	var role model.Role
	if err := s.db.Model(model.Role{}).Where("type=?", newUser.Role.Type).First(&role).Error; err != nil {
		return err
	}

	newUser.Role = model.Role{}
	newUser.RoleID = role.ID
	if err := s.db.Model(model.User{}).Create(newUser).Error; err != nil {
		return err
	}

	return nil
}

func (s *AdminRepository) ById(id uint) (*model.User, error) {
	var user model.User
	var query = s.db.Model(model.User{})
	if err := query.Where("users.id=? AND users.deleted_at IS NULL", id).
		Preload("Role").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("AdminLogs").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *AdminRepository) GetAdminLogs(req *dto.AdminLogQuery) (*utils.PaginationModel[model.AdminLog], error) {
	var query = s.db.Model(model.AdminLog{}).Joins("LEFT JOIN users ON users.id = admin_logs.user_id")
	if req != nil {
		if req.Keyword != "" {
			query = query.Where("users.email ILIKE ?", "%"+req.Keyword+"%")
			query = query.Or("users.username ILIKE ?", "%"+req.Keyword+"%")
			query = query.Or("users.display_name ILIKE ?", "%"+req.Keyword+"%")
			query = query.Or("admin_logs.action ILIKE ?", "%"+req.Keyword+"%")
			query = query.Or("admin_logs.details ILIKE ?", "%"+req.Keyword+"%")
		}
		if len(req.UserIDs) > 0 {
			query = query.Where("admin_logs.user_id IN ?", req.UserIDs)
		}
		if req.Sort != "" && req.SortBy != "" {
			query = query.Order(fmt.Sprintf("admin_logs.%s %s", req.SortBy, req.Sort))
		}
	}
	query = query.Preload("User")
	pagination, err := utils.CreatePage[model.AdminLog](query, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, err
	}
	return utils.Create(pagination, int(req.Page), int(req.Limit))
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
