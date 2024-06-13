package product

import (
	"context"

	"github.com/ariefro/buycut-api/internal/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, products []*entity.Product) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, products []*entity.Product) error {
	if err := r.db.WithContext(ctx).Create(&products).Error; err != nil {
		return err
	}

	return nil
}
