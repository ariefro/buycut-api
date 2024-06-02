// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package initializer

import (
	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/server"
)

// Injectors from initializer.go:

func InitializedServer() error {
	configConfig := config.NewLoadConfig()
	error2 := server.NewFiberServer(configConfig)
	return error2
}
