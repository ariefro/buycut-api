package company

import (
	"context"
	"fmt"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/helper"
)

type Service interface {
	Create(ctx context.Context, input CreateCompaniesRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, reqs CreateCompaniesRequest) error {
	fmt.Println("pppp", reqs)
	var companies []*entity.Company
	for _, name := range reqs.Names {
		slug := helper.GenerateSlug(name)

		company := &entity.Company{
			Name: name,
			Slug: slug,
		}

		companies = append(companies, company)
	}

	return s.repo.Create(ctx, companies)
}
