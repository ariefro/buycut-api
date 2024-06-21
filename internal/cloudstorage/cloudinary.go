package cloudstorage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadArgs struct {
	File              interface{}
	CompanyName, Slug string
	Config            *config.CloudinaryConfig
}

type UpdateArgs struct {
	Slug   string
	Config *config.CloudinaryConfig
}

type DeleteArgs struct {
	Config            *config.CloudinaryConfig
	CompanyName, Slug string
}

type UploadImageArgs struct {
	File          *multipart.FileHeader
	Company, Slug string
}

func SetupCloudinary(cfg *config.CloudinaryConfig) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryApiKey, cfg.CloudinarySecretKey)
	if err != nil {
		return nil, err
	}
	return cld, nil
}

func UploadFile(args *UploadArgs) (string, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary(args.Config)
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		PublicID: args.Slug,
		Tags:     api.CldAPIArray{args.CompanyName},
		Folder:   args.Config.CloudinaryBuycutFolder + "/" + args.CompanyName,
	}

	result, err := cld.Upload.Upload(ctx, args.File, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

func DeleteFile(args *DeleteArgs) error {
	ctx := context.Background()
	cld, err := SetupCloudinary(args.Config)
	if err != nil {
		return err
	}

	publicID := fmt.Sprintf("%s/%s/%s", args.Config.CloudinaryBuycutFolder, args.CompanyName, args.Slug)
	destroyParams := uploader.DestroyParams{PublicID: publicID}
	_, err = cld.Upload.Destroy(ctx, destroyParams)
	if err != nil {
		return err
	}

	return nil
}

func UploadImage(ctx context.Context, args *UploadImageArgs, config *config.CloudinaryConfig) (string, error) {
	if args.File == nil {
		return "", nil
	}

	if err := helper.ValidateImage(args.File); err != nil {
		return "", err
	}

	imageFile, err := args.File.Open()
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	return UploadFile(&UploadArgs{File: imageFile, Slug: args.Slug, CompanyName: args.Company, Config: config})
}
