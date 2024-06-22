package company

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/ariefro/buycut-api/config"
	cloudstorage "github.com/ariefro/buycut-api/internal/cloudstorage"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, args *createCompanyArgs) error
	Count(ctx context.Context, args *getCompaniesRequest) (int64, error)
	Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error)
	FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error)
	Update(ctx context.Context, args *updateCompanyRequest) error
	Delete(ctx context.Context, company *entity.Company) error
}

type service struct {
	config *config.Config
	repo   Repository
}

func NewService(config *config.Config, repo Repository) Service {
	return &service{config, repo}
}

type uploadImageArgs struct {
	File *multipart.FileHeader
	Slug string
}

func (s *service) Create(ctx context.Context, args *createCompanyArgs) error {
	slug := helper.GenerateSlug(args.Request.Name)

	company := &entity.Company{
		Name:        strings.ToLower(args.Request.Name),
		Slug:        slug,
		Description: args.Request.Description,
		Proof:       args.Request.Proof,
	}

	if err := s.repo.Create(ctx, company); err != nil {
		return err
	}

	imageURL, err := cloudstorage.UploadImage(ctx, &cloudstorage.UploadImageArgs{
		CompanyID: company.ID,
		File:      args.FormHeader,
		Slug:      slug,
	}, s.configureCloudinary())
	if err != nil {
		return err
	}

	if err := s.Update(ctx, &updateCompanyRequest{
		CompanyID: &company.ID,
		ImageURL:  &imageURL,
	}); err != nil {
		return err
	}

	return nil
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

func (s *service) Update(ctx context.Context, args *updateCompanyRequest) error {
	dataToUpdate := map[string]interface{}{}
	if args.Name != nil {
		slug := helper.GenerateSlug(*args.Name)
		dataToUpdate[common.ColumnName] = strings.ToLower(*args.Name)
		dataToUpdate[common.ColumnSlug] = slug
	}

	if args.Description != nil {
		dataToUpdate[common.ColumnDescription] = *args.Description
	}

	if args.ImageURL != nil {
		dataToUpdate[common.ColumnImageURL] = *args.ImageURL
	}

	return s.repo.Update(ctx, *args.CompanyID, dataToUpdate)
}

func (s *service) Delete(ctx context.Context, company *entity.Company) error {
	err := s.repo.Delete(ctx, company.ID)
	if err != nil {
		return err
	} else {
		if errDeleteFile := cloudstorage.DeleteFile(&cloudstorage.DeleteArgs{
			CompanyID: company.ID,
			Config:    s.configureCloudinary(),
			Slug:      company.Slug,
		}); errDeleteFile != nil {
			return errDeleteFile
		}
	}

	return nil
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
