//go:build wireinject
// +build wireinject

package initializer

import (
	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/server"
	"github.com/google/wire"
)

func InitializedServer() error {
	wire.Build(
		config.NewLoadConfig,
		server.NewFiberServer,
	)

	return nil
}
