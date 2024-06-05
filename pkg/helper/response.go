package helper

type BaseResponseFailed struct {
	Message string `json:"message"`
}

func ResponseFailed(message string) BaseResponseFailed {
	return BaseResponseFailed{
		Message: message,
	}
}

type BaseResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(message string, data interface{}) *BaseResponseSuccess {
	return &BaseResponseSuccess{
		Message: message,
		Data:    data,
	}
}
