package brand

import (
	"context"
	"errors"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, brands *entity.Brand) error
	FindByKeyword(ctx context.Context, keyword string) ([]*entity.Company, []*entity.Brand, error)
	FindOneByID(ctx context.Context, brandID uint) (*entity.Brand, error)
	FindAll(ctx context.Context, args *getBrandByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, []*entity.Brand, error)
	CountBrands(ctx context.Context, keyword string) (int64, error)
	Update(ctx context.Context, brandID uint, data map[string]interface{}) error
	Delete(ctx context.Context, brandID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, brands *entity.Brand) error {
	if err := r.db.WithContext(ctx).Create(&brands).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByKeyword(ctx context.Context, keyword string) ([]*entity.Company, []*entity.Brand, error) {
	var companies []*entity.Company
	var brands []*entity.Brand

	// Search in companies
	if err := r.db.WithContext(ctx).Model(&entity.Company{}).Preload("Brands").Where("LOWER(name) = LOWER(?)", keyword).Find(&companies).Error; err != nil {
		return nil, nil, err
	}

	// Search in brands
	if err := r.db.WithContext(ctx).Model(&entity.Brand{}).Preload("Company").Where("LOWER(name) = LOWER(?)", keyword).Find(&brands).Error; err != nil {
		return nil, nil, err
	}

	return companies, brands, nil
}

func (r *repository) FindOneByID(ctx context.Context, brandID uint) (*entity.Brand, error) {
	var brand *entity.Brand
	if err := r.db.WithContext(ctx).Model(&entity.Brand{}).Preload("Company").First(&brand, "id = ?", brandID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(common.BrandNotFound)
		}

		return nil, err
	}

	return brand, nil
}

func (r *repository) FindAll(ctx context.Context, args *getBrandByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, []*entity.Brand, error) {
	var companies []*entity.Company
	var brands []*entity.Brand
	keyword := "%" + args.Keyword + "%"

	// Search in companies
	resultCompanies := r.db.WithContext(ctx).Model(&entity.Company{}).Limit(paginationParams.Limit).Offset(paginationParams.Offset).Where("LOWER(name) LIKE LOWER(?)", keyword).Order("name asc").Find(&companies)
	if resultCompanies.Error != nil {
		return nil, nil, resultCompanies.Error
	}

	// Calculate limit for loading brands
	queryLimitBrand := calculateQueryLimitBrand(resultCompanies.RowsAffected, paginationParams.Limit)

	// Search in brands
	resultBrands := r.db.WithContext(ctx).Model(&entity.Brand{}).Preload("Company").Limit(int(queryLimitBrand)).Offset(paginationParams.Offset).Where("LOWER(name) LIKE LOWER(?)", keyword).Order("name asc").Find(&brands)
	if resultBrands.Error != nil {
		return nil, nil, resultBrands.Error
	}

	return companies, brands, nil
}

func (r *repository) CountBrands(ctx context.Context, keyword string) (int64, error) {
	var count int64
	key := "%" + keyword + "%"
	if err := r.db.WithContext(ctx).Model(&entity.Brand{}).Where("LOWER(name) LIKE LOWER(?)", key).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, brandID uint, data map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&entity.Brand{}).Where("id = ?", brandID).Updates(data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
			return errors.New(common.CompanyNotFound)
		}

		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(common.BrandNotFound)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, brandID uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Brand{}, "id = ?", brandID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(common.BrandNotFound)
	}

	return nil
}

// calculateQueryLimitBrand calculates the limit for loading brands based on the number of companies found
func calculateQueryLimitBrand(rowsAffected int64, limit int) int64 {
	if rowsAffected < int64(limit) {
		return int64(limit) + (int64(limit) - rowsAffected)
	}

	return int64(limit)
}
