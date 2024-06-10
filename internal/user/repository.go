package user

import (
	"context"
	"errors"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/common"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	FindOneByEmail(ctx context.Context, email string) (*entity.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (repo *repository) Create(ctx context.Context, user *entity.User) error {
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New(common.ErrDuplicateEntry)
		}

		return err
	}

	return nil
}

func (repo *repository) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := repo.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(common.EmailNotRegistered)
		}

		return nil, err
	}

	return &user, nil
}
