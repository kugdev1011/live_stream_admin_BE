package service

import (
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo  *repository.Repository
	redis *redis.Client
}

func newUserService(repo *repository.Repository, redis *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		redis: redis,
	}

}

func (s *UserService) Create(user *model.User) error {
	return s.repo.User.Create(user)
}

func (s *UserService) FindByEmail(email string) (*model.User, error) {
	return s.repo.User.FindByEmail(email)
}

func (s *UserService) FindByUsername(username string) (*model.User, error) {
	return s.repo.User.FindByUsername(username)
}
func (s *UserService) Update(user *model.User) error {
	return s.repo.User.Update(user)
}

func (s *UserService) UpdateOTP(userID uint, otp string, expiresAt time.Time) error {
	return s.repo.User.UpdateOTP(userID, otp, expiresAt)
}

func (s *UserService) ClearOTP(userID uint) error {
	return s.repo.User.ClearOTP(userID)
}
func (s *UserService) UpdatePassword(userID uint, hashedPassword string) error {
	return s.repo.User.UpdatePassword(userID, hashedPassword)
}
