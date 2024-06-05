package helper

import "github.com/gosimple/slug"

func GenerateSlug(input string) string {
	return slug.Make(input)
}
