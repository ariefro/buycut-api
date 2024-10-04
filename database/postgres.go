package database

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ariefro/buycut-api/config"
	log "github.com/sirupsen/logrus"
)

func NewConnectPostgres(config *config.Config) *gorm.DB {
	dsn := "host=" + config.PostgresHost +
		" user=" + config.PostgresUser +
		" password=" + config.PostgresPassword +
		" dbname=" + config.PostgresDatabase +
		" port=" + config.PostgresPort +
		" sslmode=disable Timezone=Asia/Jakarta"

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}))
	if err != nil {
		log.Error("cannot connect to database: ", err.Error())
		os.Exit(2)
	}

	// connection pool config
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("cannot maintain connection pools: ", err.Error())
		os.Exit(2)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	log.Info("🐝 connected successfully to the database")

	Migration(db)

	return db
}
