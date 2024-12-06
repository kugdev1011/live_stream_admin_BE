package service

import (
	"gitlab/live/be-live-api/model"
	apimodel "gitlab/live/be-live-api/model/api-model"
	"gitlab/live/be-live-api/pkg/utils"
	"gitlab/live/be-live-api/repository"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo  *repository.Repository
	redis *redis.Client
}

func (s *UserService) GetUserList(filter *apimodel.UserQuery, page, limit int) (*utils.PaginationModel[model.User], error) {
	return s.repo.User.Page(filter, page, limit)
}

func newUserService(repo *repository.Repository, redis *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		redis: redis,
	}
}
