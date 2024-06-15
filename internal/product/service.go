package product

import (
	"context"
	"errors"
	"strings"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
)

type Service interface {
	Create(ctx context.Context, args *createProductsRequest) error
	FindByKeyword(ctx context.Context, args *getProductByKeywordRequest) (interface{}, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, args *createProductsRequest) error {
	var products []*entity.Product
	for _, productName := range args.ProductNames {
		slug := helper.GenerateSlug(productName)

		product := &entity.Product{
			Name:      strings.ToLower(productName),
			Slug:      slug,
			CompanyID: args.CompanyID,
		}

		products = append(products, product)
	}

	return s.repo.Create(ctx, products)
}

func (s *service) FindByKeyword(ctx context.Context, args *getProductByKeywordRequest) (interface{}, error) {
	companies, products, err := s.repo.FindByKeyword(ctx, args.Keyword)
	if err != nil {
		return nil, err
	}

	if len(companies) > 0 {
		return companies, nil
	} else if len(products) > 0 {
		return products, nil
	} else {
		return nil, errors.New(common.ProductNotFound)
	}
}
