package service

import (
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/repository"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	User  *UserService
	Admin *AdminService
}

func NewService(repo *repository.Repository, redis *redis.Client, appConfig *conf.ApplicationConfig) *Service {
	return &Service{
		User:  newUserService(repo, redis),
		Admin: newAdminService(repo, redis, appConfig),
	}
}
