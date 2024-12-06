package service

import (
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/repository"

	"github.com/redis/go-redis/v9"
)

type AdminLogService struct {
	repo  *repository.Repository
	redis *redis.Client
}

func newAdminService(repo *repository.Repository, redis *redis.Client) *AdminLogService {
	return &AdminLogService{
		repo:  repo,
		redis: redis,
	}
}
func (s *AdminLogService) Create(adminLog *model.AdminLog) error {
	return s.repo.Admin.Create(adminLog)
}
