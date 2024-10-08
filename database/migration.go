package database

import (
	"github.com/ariefro/buycut-api/internal/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	log.Info("running migrations...")
	db.AutoMigrate(
		&entity.Company{},
		&entity.Brand{},
		&entity.User{},
	)

	log.Info("migrations complete...")
}
