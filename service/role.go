package service

import (
	"github.com/redis/go-redis/v9"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/repository"
)

type RoleService struct {
	repo  *repository.Repository
	redis *redis.Client
}

func NewRoleService(repo *repository.Repository, redis *redis.Client) *RoleService {
	return &RoleService{
		repo:  repo,
		redis: redis,
	}
}

func (s *RoleService) CreateRole(role *model.Role) error {
	return s.repo.Role.Create(role)
}

func (s *RoleService) GetRoleByID(roleID uint) (*model.Role, error) {
	return s.repo.Role.FindByID(roleID)
}

func (s *RoleService) GetRoleByType(roleType string) (*model.Role, error) {
	return s.repo.Role.FindByType(roleType)
}
