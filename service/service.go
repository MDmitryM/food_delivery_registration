package service

import (
	"context"

	models "github.com/MDmitryM/food_delivery_registration"
	"github.com/MDmitryM/food_delivery_registration/repository"
)

type Servicer interface {
	CreateUser(ctx context.Context, arg repository.CreateUserParams) (int, error)
	DeleteUserByID(ctx context.Context, id int32) (int64, error)
	GetUserByID(ctx context.Context, id int32) (repository.User, error)
	UpdateUserPwd(ctx context.Context, arg repository.UpdateUserPwdParams) (repository.User, error)
}

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, arg models.User) (int32, error) {
	return s.repo.CreateUser(ctx, repository.CreateUserParams{Login: arg.Login, PwdHash: arg.Password})
}

func (s *Service) DeleteUserByID(ctx context.Context, id int32) (int64, error) {
	return s.repo.DeleteUserByID(ctx, id)
}

func (s *Service) GetUserByID(ctx context.Context, id int32) (repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) UpdateUserPwd(ctx context.Context, arg models.UpdateUser) (repository.User, error) {
	return s.repo.UpdateUserPwd(ctx, repository.UpdateUserPwdParams{ID: arg.ID, PwdHash: arg.PwdHash})
}

//TODO: add JWT bisness logic
