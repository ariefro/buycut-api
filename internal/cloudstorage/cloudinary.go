package cloudstorage

import (
	"context"

	"github.com/ariefro/buycut-api/config"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2"
)

type UploadArgs struct {
	File   interface{}
	Slug   string
	Config *config.CloudinaryConfig
}

func setupCloudinary(config *config.CloudinaryConfig) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.SecretKey)
	if err != nil {
		return nil, err
	}

	return cld, nil
}

func Upload(publicID string, args *UploadArgs) (string, error) {
	ctx := context.Background()
	cld, err := setupCloudinary(args.Config)
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		PublicID: publicID,
		Tags:     api.CldAPIArray{args.Slug},
		Folder:   args.Config.BuycutFolder + "/" + args.Slug,
	}

	result, err := cld.Upload.Upload(ctx, args.File, uploadParams)
	if err != nil {
		return "", err
	}

	imageURL := result.SecureURL
	return imageURL, nil
}
