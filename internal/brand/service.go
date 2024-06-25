package brand

import (
	"context"
	"errors"
	"sort"

	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/cloudstorage"
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, args *createBrandArgs) error
	FindByKeyword(ctx context.Context, args *getBrandByKeywordRequest) (interface{}, error)
	FindOneByID(ctx context.Context, brandID uint) (*entity.Brand, error)
	FindAll(ctx context.Context, args *getBrandByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*boycottedResult, error)
	CountAll(ctx context.Context, args *getBrandByKeywordRequest) (int64, error)
	Update(ctx context.Context, brandID uint, args *updateBrandArgs) error
	Delete(ctx context.Context, brand *entity.Brand) error
}

type service struct {
	config      *config.Config
	repo        Repository
	companyRepo company.Repository
}

func NewService(config *config.Config, repo Repository, companyRepo company.Repository) Service {
	return &service{config, repo, companyRepo}
}

func (s *service) Create(ctx context.Context, args *createBrandArgs) error {
	slug := helper.GenerateSlug(args.Request.Name)
	imageURL, err := cloudstorage.UploadImage(ctx, &cloudstorage.UploadImageArgs{
		CompanyID: args.CompanyID,
		File:      args.FormHeader,
		Slug:      slug,
	}, s.configureCloudinary())
	if err != nil {
		return err
	}

	brand := &entity.Brand{
		Name:      args.Request.Name,
		Slug:      slug,
		CompanyID: args.Request.CompanyID,
		ImageURL:  imageURL,
	}

	return s.repo.Create(ctx, brand)
}

func (s *service) FindOneByID(ctx context.Context, brandID uint) (*entity.Brand, error) {
	return s.repo.FindOneByID(ctx, brandID)
}

func (s *service) FindByKeyword(ctx context.Context, args *getBrandByKeywordRequest) (interface{}, error) {
	companies, brands, err := s.repo.FindByKeyword(ctx, args.Keyword)
	if err != nil {
		return nil, err
	}

	if len(companies) > 0 {
		return companies, nil
	} else if len(brands) > 0 {
		return brands, nil
	} else {
		return nil, errors.New(common.BrandNotFound)
	}
}

func (s *service) FindAll(ctx context.Context, args *getBrandByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*boycottedResult, error) {
	companies, brands, err := s.repo.FindAll(ctx, args, paginationParams)
	if err != nil {
		return nil, err
	}

	var results []*boycottedResult
	for _, company := range companies {
		results = append(results, &boycottedResult{
			ID:          company.ID,
			Name:        company.Name,
			Slug:        company.Slug,
			Description: company.Description,
			ImageURL:    company.ImageURL,
			Proof:       company.Proof,
			Company:     nil,
			Type:        "company",
		})
	}

	for _, brand := range brands {
		results = append(results, &boycottedResult{
			ID:          brand.ID,
			Name:        brand.Name,
			Slug:        brand.Slug,
			Description: brand.Company.Description,
			ImageURL:    brand.ImageURL,
			Proof:       brand.Company.Proof,
			Company:     brand.Company,
			Type:        "brand",
		})
	}

	// Sort results by name in ascending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func (s *service) CountAll(ctx context.Context, args *getBrandByKeywordRequest) (int64, error) {
	companyCount, err := s.companyRepo.CountCompanies(ctx, args.Keyword)
	if err != nil {
		return 0, err
	}

	brandCount, err := s.repo.CountBrands(ctx, args.Keyword)
	if err != nil {
		return 0, err
	}

	results := companyCount + brandCount

	return results, nil
}

func (s *service) Update(ctx context.Context, brandID uint, args *updateBrandArgs) error {
	dataToUpdate := map[string]interface{}{}

	slug := helper.GenerateSlug(args.Request.Name)
	dataToUpdate[common.ColumnName] = args.Request.Name
	dataToUpdate[common.ColumnSlug] = slug

	if args.Request.CompanyID != nil {
		dataToUpdate[common.ColumnCompanyID] = args.Request.CompanyID
	}

	if args.FormHeader != nil {
		// jika nama dari produk tidak sama dengan nama dari request, maka hapus file yang lama
		if args.Brand.Name != args.Request.Name {
			cloudstorage.DeleteFile(&cloudstorage.DeleteArgs{
				CompanyID: args.Brand.Company.ID,
				Config:    s.configureCloudinary(),
				Slug:      args.Brand.Slug,
			})
		}

		imageURL, err := cloudstorage.UploadImage(ctx, &cloudstorage.UploadImageArgs{
			CompanyID: args.Brand.Company.ID,
			File:      args.FormHeader,
			Slug:      slug,
		}, s.configureCloudinary())
		if err != nil {
			return err
		}

		dataToUpdate[common.ColumnImageURL] = imageURL
	}

	return s.repo.Update(ctx, brandID, dataToUpdate)
}

func (s *service) Delete(ctx context.Context, brand *entity.Brand) error {
	err := s.repo.Delete(ctx, brand.ID)
	if err != nil {
		return err
	} else {
		if errDeleteFile := cloudstorage.DeleteFile(&cloudstorage.DeleteArgs{
			CompanyID: brand.Company.ID,
			Config:    s.configureCloudinary(),
			Slug:      brand.Slug,
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
