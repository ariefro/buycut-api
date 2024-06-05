//go:build wireinject
// +build wireinject

package initializer

import (
	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/database"
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/server"
	"github.com/google/wire"
)

var companySet = wire.NewSet(
	company.NewRepository,
	company.NewService,
	company.NewController,
)

func InitializedServer() error {
	wire.Build(
		config.NewLoadConfig,
		database.NewConnectPostgres,
		companySet,
		server.NewFiberServer,
	)

	return nil
}
