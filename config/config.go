package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort       string `mapstructure:"APP_PORT"`
	ClientBaseURL string `mapstructure:"CLIENT_BASE_URL"`

	CloudinaryApiKey       string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryBuycutFolder string `mapstructure:"CLOUDINARY_BUYCUT_FOLDER"`
	CloudinaryCloudName    string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinarySecretKey    string `mapstructure:"CLOUDINARY_SECRET_KEY"`

	JwtAccessTokenSecret   string `mapstructure:"JWT_SECRET_KEY"`
	JwtAccessTokenDuration uint   `mapstructure:"JWT_ACCESS_TOKEN_DURATION"`

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
