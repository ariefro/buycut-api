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
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (repo *repository) Create(ctx context.Context, companies []*entity.Company) error {
	if err := repo.db.WithContext(ctx).Create(companies).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New(common.ErrDuplicateEntry)
		}

		return err
	}

	return nil
}
