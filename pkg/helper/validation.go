package helper

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/ariefro/buycut-api/pkg/common"
)

func ValidateImage(file *multipart.FileHeader) error {
	// Check for supported image extensions
	if err := validateImageExtension(file.Filename); err != nil {

		return err
	}

	// Check for maximum file size
	if err := validateFileSize(file.Size); err != nil {
		return err
	}

	return nil
}

func validateImageExtension(filename string) error {
	validExtensions := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if _, ok := validExtensions[ext]; !ok {
		return errors.New(common.InvalidImageFile)
	}

	return nil
}

func validateFileSize(size int64) error {
	if size > 1*1024*1024 { // 1 MB limit
		return errors.New(common.FileSizeIsTooLarge)
	}

	return nil
}
