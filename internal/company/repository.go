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
	Create(ctx context.Context, companies []*entity.Company) error
	Count(ctx context.Context, args *getCompaniesRequest) (int64, error)
	Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error)
	FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error)
	Update(ctx context.Context, companyID uint, data map[string]interface{}) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, companies []*entity.Company) error {
	if err := r.db.WithContext(ctx).Create(companies).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New(common.ErrDuplicateEntry)
		}

		return err
	}

	return nil
}

func (r *repository) Count(ctx context.Context, args *getCompaniesRequest) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Company{})

	if args.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+args.Keyword+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count companies: %w", err)
	}

	return count, nil
}

func (r *repository) Find(ctx context.Context, args *getCompaniesRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, error) {
	var companies []*entity.Company
	query := r.db.WithContext(ctx).Model(&entity.Company{})

	if args.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+args.Keyword+"%")
	}

	if err := query.Limit(paginationParams.Limit).Offset(paginationParams.Offset).Find(&companies).Error; err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *repository) FindOneByID(ctx context.Context, companyID uint) (*entity.Company, error) {
	var company entity.Company
	if err := r.db.WithContext(ctx).Model(&entity.Company{}).First("id = ?", companyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(common.CompanyNotFound)
		}

		return nil, err
	}

	return &company, nil
}

func (r *repository) Update(ctx context.Context, companyID uint, data map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(&entity.Company{}).Where("id = ?", companyID).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
