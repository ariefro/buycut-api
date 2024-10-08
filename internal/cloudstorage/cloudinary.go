package cloudstorage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadArgs struct {
	File      interface{}
	CompanyID uint
	Slug      string
	Config    *config.CloudinaryConfig
}

type UpdateArgs struct {
	Slug   string
	Config *config.CloudinaryConfig
}

type DeleteArgs struct {
	CompanyID uint
	Config    *config.CloudinaryConfig
	Slug      string
}

type DeleteAssetsByTagArgs struct {
	CompanyID uint
	Config    *config.CloudinaryConfig
}

type DeleteEmptyFolderArgs struct {
	CompanyID uint
	Config    *config.CloudinaryConfig
}

type UploadImageArgs struct {
	CompanyID uint
	File      *multipart.FileHeader
	Slug      string
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

	companyIDStr := strconv.FormatUint(uint64(args.CompanyID), 10)
	uploadParams := uploader.UploadParams{
		PublicID: args.Slug,
		Tags:     api.CldAPIArray{companyIDStr},
		Folder:   args.Config.CloudinaryBuycutFolder + "/" + companyIDStr,
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

	companyIDStr := strconv.FormatUint(uint64(args.CompanyID), 10)
	publicID := fmt.Sprintf("%s/%s/%s", args.Config.CloudinaryBuycutFolder, companyIDStr, args.Slug)
	destroyParams := uploader.DestroyParams{PublicID: publicID}
	_, err = cld.Upload.Destroy(ctx, destroyParams)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAssetsByTag(args *DeleteAssetsByTagArgs) error {
	ctx := context.Background()
	cld, err := SetupCloudinary(args.Config)
	if err != nil {
		return err
	}

	companyIDStr := strconv.FormatUint(uint64(args.CompanyID), 10)
	_, err = cld.Admin.DeleteAssetsByTag(ctx, admin.DeleteAssetsByTagParams{Tag: companyIDStr})
	if err != nil {
		return err
	}

	return nil
}

func DeleteEmptyFolder(args *DeleteEmptyFolderArgs) error {
	ctx := context.Background()
	cld, err := SetupCloudinary(args.Config)
	if err != nil {
		return err
	}

	companyIDStr := strconv.FormatUint(uint64(args.CompanyID), 10)
	pathFolder := fmt.Sprintf("%s/%s", args.Config.CloudinaryBuycutFolder, companyIDStr)
	_, err = cld.Admin.DeleteFolder(ctx, admin.DeleteFolderParams{Folder: pathFolder})
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

	return UploadFile(&UploadArgs{File: imageFile, Slug: args.Slug, CompanyID: args.CompanyID, Config: config})
}
