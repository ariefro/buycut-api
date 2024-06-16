package cloudinary

import (
	"github.com/ariefro/buycut-api/config"
	cld "github.com/cloudinary/cloudinary-go/v2"
)

func setupCloudinary(config *config.CloudinaryConfig) (*cld.Cloudinary, error) {
	cld, err := cld.NewFromParams(config.CloudName, config.APIKey, config.SecretKey)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
