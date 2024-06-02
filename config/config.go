package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort       string `mapstructure:"APP_PORT"`
	ClientBaseURL string `mapstructure:"CLIENT_BASE_URL"`

	PostgresDatabase string `mapstructure:"POSTGRES_DATABASE"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
}

func NewLoadConfig() *Config {
	log.Println("configuring server...")
	env := os.Getenv("APP_ENV")

	viper.AddConfigPath(".")
	viper.SetConfigName(".env." + env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var config *Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return config
}
