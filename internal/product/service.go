package product

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/cloudstorage"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
)

type Service interface {
	Create(ctx context.Context, args *createProductsRequest, formHeader *multipart.FileHeader) error
	FindByKeyword(ctx context.Context, args *getProductByKeywordRequest) (interface{}, error)
}

type service struct {
	config *config.Config
	repo   Repository
}

func NewService(config *config.Config, repo Repository) Service {
	return &service{config, repo}
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

func (s *service) configureCloudinary() *config.CloudinaryConfig {
	var config = &config.CloudinaryConfig{
		CloudinaryCloudName:    s.config.CloudinaryCloudName,
		CloudinaryApiKey:       s.config.CloudinaryApiKey,
		CloudinarySecretKey:    s.config.CloudinarySecretKey,
		CloudinaryBuycutFolder: s.config.CloudinaryBuycutFolder,
	}

	return config
}
