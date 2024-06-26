package company

import (
	"context"
	"errors"
	"fmt"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, companies *entity.Company) error
	CountCompanies(ctx context.Context, keyword string) (int64, error)
	Count(ctx context.Context) (int64, error)
	Find(ctx context.Context, paginationParams *pagination.PaginationParams) ([]*entity.Company, error)
	FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error)
	Update(ctx context.Context, companyID uint, data map[string]interface{}) error
	DeleteAssociateCompanyBrandsInTx(ctx context.Context, tx *gorm.DB, companyID uint) error
	DeleteInTx(ctx context.Context, tx *gorm.DB, companyID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, company *entity.Company) error {
	if err := r.db.WithContext(ctx).Create(company).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) CountCompanies(ctx context.Context, keyword string) (int64, error) {
	var count int64
	key := "%" + keyword + "%"
	if err := r.db.WithContext(ctx).Model(&entity.Company{}).Where("LOWER(name) LIKE LOWER(?)", key).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Company{})

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count companies: %w", err)
	}

	return count, nil
}

func (r *repository) Find(ctx context.Context, paginationParams *pagination.PaginationParams) ([]*entity.Company, error) {
	var companies []*entity.Company
	query := r.db.WithContext(ctx).Model(&entity.Company{})

	if err := query.Limit(paginationParams.Limit).Offset(paginationParams.Offset).Order("name asc").Find(&companies).Error; err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *repository) FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error) {
	var company *entity.Company
	if err := r.db.WithContext(ctx).Preload("Brands").First(&company, "id = ?", companyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(common.CompanyNotFound)
		}

		return nil, err
	}

	return company, nil
}

func (r *repository) Update(ctx context.Context, companyID uint, data map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&entity.Company{}).Where("id = ?", companyID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(common.CompanyNotFound)
	}

	return nil
}

func (r *repository) DeleteAssociateCompanyBrandsInTx(ctx context.Context, tx *gorm.DB, companyID uint) error {
	if err := tx.WithContext(ctx).Delete(&entity.Brand{}, "company_id = ?", companyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(common.CompanyNotFound)
		}

		return err
	}

	return nil
}

func (r *repository) DeleteInTx(ctx context.Context, tx *gorm.DB, companyID uint) error {
	result := tx.WithContext(ctx).Delete(&entity.Company{}, "id = ?", companyID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(common.CompanyNotFound)
	}

	return nil
}
