package company

import (
	"context"
	"mime/multipart"

	"github.com/ariefro/buycut-api/config"
	cloudstorage "github.com/ariefro/buycut-api/internal/cloudstorage"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, args *createCompanyArgs) error
	Count(ctx context.Context, args *getCompaniesRequest) (int64, error)
	Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error)
	FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error)
	Update(ctx context.Context, args *updateCompanyArgs) error
	Delete(ctx context.Context, company *entity.Company) error
}

type service struct {
	db     *gorm.DB
	config *config.Config
	repo   Repository
}

func NewService(db *gorm.DB, config *config.Config, repo Repository) Service {
	return &service{db, config, repo}
}

type uploadImageArgs struct {
	File *multipart.FileHeader
	Slug string
}

func (s *service) Create(ctx context.Context, args *createCompanyArgs) error {
	slug := helper.GenerateSlug(args.Request.Name)
	company := &entity.Company{
		Name:        args.Request.Name,
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

	if err := s.Update(ctx, &updateCompanyArgs{
		Request: &updateCompanyRequest{
			CompanyID: company.ID,
			ImageURL:  &imageURL,
		},
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

func (s *service) Update(ctx context.Context, args *updateCompanyArgs) error {
	dataToUpdate := map[string]interface{}{}
	var slug string
	if args.Request.Name != nil {
		slug = helper.GenerateSlug(*args.Request.Name)
		dataToUpdate[common.ColumnName] = *args.Request.Name
		dataToUpdate[common.ColumnSlug] = slug
	}

	if args.Request.Description != nil {
		dataToUpdate[common.ColumnDescription] = *args.Request.Description
	}

	if args.Request.Proof != nil {
		dataToUpdate[common.ColumnProof] = pq.StringArray(args.Request.Proof)
	}

	if args.FormHeader != nil {
		// jika tidak ada inputan nama, set slug dari current company
		if args.Request.Name == nil {
			slug = args.Company.Slug
		}

		// jika ada inputan nama dan tidak sama dengan nama dari current company, hapus file lama di cloud
		if args.Request.Name != nil && *args.Request.Name != args.Company.Name {
			if err := cloudstorage.DeleteFile(&cloudstorage.DeleteArgs{
				CompanyID: args.Company.ID,
				Config:    s.configureCloudinary(),
				Slug:      args.Company.Slug,
			}); err != nil {
				return err
			}
		}

		imageURL, err := cloudstorage.UploadImage(ctx, &cloudstorage.UploadImageArgs{
			CompanyID: args.Request.CompanyID,
			File:      args.FormHeader,
			Slug:      slug,
		}, s.configureCloudinary())
		if err != nil {
			return err
		}

		dataToUpdate[common.ColumnImageURL] = imageURL
	}

	if args.Request.ImageURL != nil {
		dataToUpdate[common.ColumnImageURL] = *args.Request.ImageURL
	}

	return s.repo.Update(ctx, args.Request.CompanyID, dataToUpdate)
}

func (s *service) Delete(ctx context.Context, company *entity.Company) error {
	if errTx := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteAssociateCompanyBrandsInTx(ctx, tx, company.ID); err != nil {
			return err
		}

		if err := s.repo.DeleteInTx(ctx, tx, company.ID); err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return errTx
	}

	if err := cloudstorage.DeleteAssetsByTag(&cloudstorage.DeleteAssetsByTagArgs{
		CompanyID: company.ID,
		Config:    s.configureCloudinary(),
	}); err != nil {
		return err
	}

	if err := cloudstorage.DeleteEmptyFolder(&cloudstorage.DeleteEmptyFolderArgs{
		CompanyID: company.ID,
		Config:    s.configureCloudinary(),
	}); err != nil {
		return err
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
