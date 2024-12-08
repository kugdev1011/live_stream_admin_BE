package service

import (
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/model"
	apimodel "gitlab/live/be-live-api/model/api-model"
	"gitlab/live/be-live-api/pkg/utils"
	"gitlab/live/be-live-api/repository"

	"github.com/redis/go-redis/v9"
)

type AdminService struct {
	repo      *repository.Repository
	redis     *redis.Client
	appConfig *conf.ApplicationConfig
}

func (s *AdminService) CreateAdmin(request *apimodel.CreateAdminRequest) (*model.User, error) {
	var newUser = new(model.User)
	newUser.Username = request.UserName
	newUser.PasswordHash = utils.HashingPassword(request.Password, s.appConfig.SaltKey)
	newUser.DisplayName = request.DisplayName
	newUser.Email = request.Email
	newUser.CreatedByID = request.CreatedByID
	newUser.UpdatedByID = request.CreatedByID

	return s.repo.Admin.CreateAdmin(newUser)
}

func (s *AdminService) ById(id uint) (*model.User, error) {
	return s.repo.Admin.ById(id)
}

func newAdminService(repo *repository.Repository, redis *redis.Client, appConfig *conf.ApplicationConfig) *AdminService {
	return &AdminService{
		repo:      repo,
		redis:     redis,
		appConfig: appConfig,
	}
}
