package user

import (
	"context"
	"errors"

	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
)

type Service interface {
	Register(ctx context.Context, req *registerRequest) error
	Login(ctx context.Context, req *loginRequest, user *entity.User) (string, error)
	FindOneByEmail(ctx context.Context, email string) (*entity.User, error)
}

type service struct {
	config *config.Config
	repo   Repository
}

func NewService(config *config.Config, repo Repository) Service {
	return &service{config, repo}
}

func (s *service) Register(ctx context.Context, req *registerRequest) error {
	hashedPassword, err := helper.HashedPassword(req.Password)
	if err != nil {
		return err
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.repo.Create(ctx, &user)
}

func (s *service) Login(ctx context.Context, req *loginRequest, user *entity.User) (string, error) {
	if err := helper.CheckPassword(req.Password, user.Password); err != nil {
		return "", errors.New(common.ErrInvalidEmailOrPassword)
	}

	tokenArgs := helper.GenerateAccessTokenArgs{
		UserID:        user.ID,
		TokenDuration: s.config.JwtAccessTokenDuration,
		SecretKey:     s.config.JwtAccessTokenSecret,
	}

	accessToken, err := helper.GenerateAccessToken(&tokenArgs)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *service) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.FindOneByEmail(ctx, email)
}
