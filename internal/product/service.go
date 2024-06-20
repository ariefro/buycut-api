package product

import (
	"context"
	"errors"
	"mime/multipart"
	"sort"
	"strings"

	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/cloudstorage"
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, args *createProductsRequest, formHeader *multipart.FileHeader) error
	FindByKeyword(ctx context.Context, args *getProductByKeywordRequest) (interface{}, error)
	FindAll(ctx context.Context, args *getProductByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*boycottedResult, error)
	CountAll(ctx context.Context, args *getProductByKeywordRequest) (int64, error)
}

type service struct {
	config      *config.Config
	repo        Repository
	companyRepo company.Repository
}

func NewService(config *config.Config, repo Repository, companyRepo company.Repository) Service {
	return &service{config, repo, companyRepo}
}

func (s *service) Create(ctx context.Context, args *createProductsRequest, formHeader *multipart.FileHeader) error {
	slug := helper.GenerateSlug(args.Name)
	imageURL, err := cloudstorage.UploadImage(ctx, &cloudstorage.UploadImageArgs{
		File: formHeader,
		Slug: slug,
	}, s.configureCloudinary())
	if err != nil {
		return err
	}

	product := &entity.Product{
		Name:      strings.ToLower(args.Name),
		Slug:      slug,
		CompanyID: args.CompanyID,
		ImageURL:  imageURL,
	}

	return s.repo.Create(ctx, product)
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

func (s *service) FindAll(ctx context.Context, args *getProductByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*boycottedResult, error) {
	companies, products, err := s.repo.FindAll(ctx, args, paginationParams)
	if err != nil {
		return nil, err
	}

	var results []*boycottedResult
	for _, company := range companies {
		results = append(results, &boycottedResult{
			ID:          company.ID,
			Name:        helper.MakeTitle(company.Name),
			Slug:        company.Slug,
			Description: company.Description,
			ImageURL:    company.ImageURL,
			Type:        "company",
		})
	}

	for _, product := range products {
		results = append(results, &boycottedResult{
			ID:          product.ID,
			Name:        helper.MakeTitle(product.Name),
			Slug:        product.Slug,
			Description: product.Company.Description,
			ImageURL:    product.ImageURL,
			Type:        "product",
		})
	}

	// Sort results by name in ascending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func (s *service) CountAll(ctx context.Context, args *getProductByKeywordRequest) (int64, error) {
	companyCount, err := s.companyRepo.CountCompanies(ctx, args.Keyword)
	if err != nil {
		return 0, err
	}

	productCount, err := s.repo.CountProducts(ctx, args.Keyword)
	if err != nil {
		return 0, err
	}

	results := companyCount + productCount

	return results, nil
}

func (s *service) configureCloudinary() *config.CloudinaryConfig {
	var config = &config.CloudinaryConfig{
		CloudinaryCloudName:    s.config.CloudinaryCloudName,
		CloudinaryApiKey:       s.config.CloudinaryApiKey,
		CloudinarySecretKey:    s.config.CloudinarySecretKey,
		CloudinaryBuycutFolder: s.config.CloudinaryBuycutFolder,
	}

	return config
}
