package service

import (
	"gitlab/live/be-live-api/repository"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	User  *UserService
	admin *AdminService
}

func NewService(repo *repository.Repository, redis *redis.Client) *Service {
	return &Service{
		User:  newUserService(repo, redis),
		admin: newAdminService(repo, redis),
	}
}
