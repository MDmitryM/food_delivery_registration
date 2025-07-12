package service

import (
	"context"

	models "github.com/MDmitryM/food_delivery_registration"
	"github.com/MDmitryM/food_delivery_registration/repository"
	"golang.org/x/crypto/bcrypt"
)

type Servicer interface {
	CreateUser(ctx context.Context, arg repository.CreateUserParams) (int, string, error)
	DeleteUserByID(ctx context.Context, id int32) (int64, error)
	GetUserByID(ctx context.Context, id int32) (repository.User, error)
	UpdateUserPwd(ctx context.Context, arg repository.UpdateUserPwdParams) (repository.User, error)
	IsUserValid(ctx context.Context, arg models.User) (string, error)
}

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, arg models.User) (int32, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}

	arg.Password = string(hash)

	userID, err := s.repo.CreateUser(ctx, repository.CreateUserParams{Login: arg.Login, PwdHash: arg.Password})
	if err != nil {
		return 0, "", err
	}

	token, err := s.GenerateToken(userID)
	if err != nil {
		return 0, "", err
	}

	return userID, token, nil
}

func (s *Service) DeleteUserByID(ctx context.Context, id int32) (int64, error) {
	return s.repo.DeleteUserByID(ctx, id)
}

func (s *Service) GetUserByID(ctx context.Context, id int32) (repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) UpdateUserPwd(ctx context.Context, arg models.UpdateUser) (repository.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return repository.User{}, err
	}
	arg.Password = string(hash)

	return s.repo.UpdateUserPwd(ctx, repository.UpdateUserPwdParams{ID: arg.ID, PwdHash: arg.Password})
}

func (s *Service) IsUserValid(ctx context.Context, arg models.User) (string, error) {
	user, err := s.repo.IsUserValid(ctx, arg.Login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PwdHash), []byte(arg.Password)); err != nil {
		return "", err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

//TODO: add JWT bisness logic
