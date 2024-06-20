package product

import (
	"context"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, products *entity.Product) error
	FindByKeyword(ctx context.Context, keyword string) ([]*entity.Company, []*entity.Product, error)
	FindAll(ctx context.Context, args *getProductByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, []*entity.Product, error)
	CountProducts(ctx context.Context, keyword string) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, products *entity.Product) error {
	if err := r.db.WithContext(ctx).Create(&products).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByKeyword(ctx context.Context, keyword string) ([]*entity.Company, []*entity.Product, error) {
	var companies []*entity.Company
	var products []*entity.Product

	// Search in companies
	if err := r.db.WithContext(ctx).Model(&entity.Company{}).Preload("Products").Where("LOWER(name) = LOWER(?)", keyword).Find(&companies).Error; err != nil {
		return nil, nil, err
	}

	// Search in products
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Preload("Company").Where("LOWER(name) = LOWER(?)", keyword).Find(&products).Error; err != nil {
		return nil, nil, err
	}

	return companies, products, nil
}

func (r *repository) FindAll(ctx context.Context, args *getProductByKeywordRequest, paginationParams *pagination.PaginationParams) ([]*entity.Company, []*entity.Product, error) {
	var companies []*entity.Company
	var products []*entity.Product
	keyword := "%" + args.Keyword + "%"

	// Search in companies
	if err := r.db.WithContext(ctx).Model(&entity.Company{}).Limit(paginationParams.Limit).Offset(paginationParams.Offset).Where("LOWER(name) LIKE LOWER(?)", keyword).Find(&companies).Error; err != nil {
		return nil, nil, err
	}

	// Search in products
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Preload("Company").Limit(paginationParams.Limit).Offset(paginationParams.Offset).Where("LOWER(name) LIKE LOWER(?)", keyword).Find(&products).Error; err != nil {
		return nil, nil, err
	}

	return companies, products, nil
}

func (r *repository) CountProducts(ctx context.Context, keyword string) (int64, error) {
	var count int64
	key := "%" + keyword + "%"
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("LOWER(name) LIKE LOWER(?)", key).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
