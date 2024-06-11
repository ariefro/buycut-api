package company

import (
	"context"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, args createCompaniesRequest) error
	Count(ctx context.Context, args *getCompaniesRequest) (int64, error)
	Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error)
	FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error)
	Update(ctx context.Context, args *updateCompaniesRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, reqs createCompaniesRequest) error {
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

func (s *service) Count(ctx context.Context, args *getCompaniesRequest) (int64, error) {
	return s.repo.Count(ctx, args)
}

func (s *service) Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error) {
	return s.repo.Find(ctx, args, paginationParams)
}

func (s *service) FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error) {
	return s.repo.FindOneByID(ctx, companyID)
}

func (s *service) Update(ctx context.Context, args *updateCompaniesRequest) error {
	slug := helper.GenerateSlug(args.Name)
	dataToUpdate := map[string]interface{}{
		common.ColumnName: args.Name,
		common.ColumnSlug: slug,
	}

	return s.repo.Update(ctx, args.CompanyID, dataToUpdate)
}
