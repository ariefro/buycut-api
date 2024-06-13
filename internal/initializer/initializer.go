//go:build wireinject
// +build wireinject

package initializer

import (
	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/database"
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/product"
	"github.com/ariefro/buycut-api/internal/server"
	"github.com/ariefro/buycut-api/internal/user"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	user.NewRepository,
	user.NewService,
	user.NewController,
)

var companySet = wire.NewSet(
	company.NewRepository,
	company.NewService,
	company.NewController,
)

var productSet = wire.NewSet(
	product.NewRepository,
	product.NewService,
	product.NewController,
)

func InitializedServer() error {
	wire.Build(
		config.NewLoadConfig,
		database.NewConnectPostgres,
		userSet,
		companySet,
		productSet,
		server.NewFiberServer,
	)

	return nil
}
