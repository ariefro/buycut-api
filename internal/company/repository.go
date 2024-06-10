package company

import (
	"context"
	"errors"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, companies []*entity.Company) error
	Find(ctx context.Context, keyword string) ([]*entity.Company, error)
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

func (r *repository) Find(ctx context.Context, keyword string) ([]*entity.Company, error) {
	var companies []*entity.Company

	query := r.db.WithContext(ctx).Model(&entity.Company{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Find(&companies).Error; err != nil {
		return nil, err
	}

	return companies, nil
}
