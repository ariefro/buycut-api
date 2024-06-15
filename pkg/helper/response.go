package helper

import "github.com/ariefro/buycut-api/pkg/pagination"

type baseResponseFailed struct {
	Message string `json:"message"`
}

func ResponseFailed(message string) baseResponseFailed {
	return baseResponseFailed{
		Message: message,
	}
}

type baseResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(message string, data interface{}) *baseResponseSuccess {
	return &baseResponseSuccess{
		Message: message,
		Data:    data,
	}
}

type baseResponseSuccessWithPagination struct {
	Message string            `json:"message"`
	Pages   *pagination.Pages `json:"page"`
	Data    interface{}       `json:"data"`
}

func ResponseSuccessWithPagination(message string, data interface{}, pages *pagination.Pages) baseResponseSuccessWithPagination {
	return baseResponseSuccessWithPagination{
		Message: message,
		Pages:   pages,
		Data:    data,
	}
}
