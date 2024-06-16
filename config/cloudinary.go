package config

type CloudinaryConfig struct {
	CloudName    string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	APIKey       string `mapstructure:"CLOUDINARY_API_KEY"`
	SecretKey    string `mapstructure:"CLOUDINARY_SECRET_KEY"`
	BuycutFolder string `mapstructure:"CLOUDINARY_BUYCUT_FOLDER"`
}
