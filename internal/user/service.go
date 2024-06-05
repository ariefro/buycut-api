package user

import (
	"context"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/helper"
)

type Service interface {
	Register(ctx context.Context, req *registerRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (svc *service) Register(ctx context.Context, req *registerRequest) error {
	hashedPassword, err := helper.HashedPassword(req.Password)
	if err != nil {
		return err
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	return svc.repo.Create(ctx, &user)
}
