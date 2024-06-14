package product

import (
	"context"
	"strings"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/helper"
)

type Service interface {
	Create(ctx context.Context, args *createProductsRequest) error
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
