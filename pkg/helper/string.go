package helper

import (
	"github.com/gosimple/slug"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateSlug(input string) string {
	return slug.Make(input)
}

func MakeTitle(input string) string {
	return cases.Title(language.English).String(input)
}
