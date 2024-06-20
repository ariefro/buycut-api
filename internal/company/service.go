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
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, args *createCompaniesRequest, formHeader *multipart.FileHeader) error
	Count(ctx context.Context, args *getCompaniesRequest) (int64, error)
	Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error)
	FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error)
	Update(ctx context.Context, args *updateCompaniesRequest) error
	Delete(ctx context.Context, companyID uint) error
}

type service struct {
	config *config.Config
	repo   Repository
}

func NewService(config *config.Config, repo Repository) Service {
	return &service{config, repo}
}

type uploadImageArgs struct {
	File     *multipart.FileHeader
	Slug     string
	PublicID string
}

func (s *service) Create(ctx context.Context, args *createCompaniesRequest, formHeader *multipart.FileHeader) error {
	publicID := uuid.NewString()
	slug := helper.GenerateSlug(args.Name)

	imageURL, err := s.uploadImage(ctx, &uploadImageArgs{
		File:     formHeader,
		Slug:     slug,
		PublicID: publicID,
	})
	if err != nil {
		return err
	}

	company := &entity.Company{
		Name:        strings.ToLower(args.Name),
		Slug:        slug,
		Description: args.Description,
		ImageURL:    imageURL,
	}

	return s.repo.Create(ctx, company)
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
		common.ColumnName: strings.ToLower(args.Name),
		common.ColumnSlug: slug,
	}

	return s.repo.Update(ctx, args.CompanyID, dataToUpdate)
}

func (s *service) Delete(ctx context.Context, companyID uint) error {
	return s.repo.Delete(ctx, companyID)
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

func (s *service) uploadImage(ctx context.Context, args *uploadImageArgs) (string, error) {
	if args.File == nil {
		return "", nil
	}

	err := helper.ValidateImage(args.File)
	if err != nil {
		return "", err
	}

	imageFile, err := args.File.Open()
	if err != nil {
		return "", err
	}
	defer imageFile.Close() // Ensure file is closed even on errors

	cloudinaryConfig := s.configureCloudinary()
	imageURL, err := cloudstorage.Upload(&cloudstorage.UploadArgs{File: imageFile, Slug: args.Slug, Config: cloudinaryConfig})
	if err != nil {
		return "", err
	}

	return imageURL, nil
}
