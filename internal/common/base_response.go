package common

type MapResponse = BaseResponse[map[string]interface{}]

// BaseResponse 通用响应类
type BaseResponse[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
}

// NewBaseResponse 构造响应
func NewBaseResponse[T any](code int, data T, message string) BaseResponse[T] {
	return BaseResponse[T]{Code: code, Data: data, Message: message}
}

// NewBaseResponseFromErrorCode 从错误码构造响应（无 data）
func NewBaseResponseFromErrorCode(code int, message string) BaseResponse[any] {
	return BaseResponse[any]{Code: code, Message: message}
}
