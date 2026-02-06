package common

import "aicode/internal/exception"

// Success 成功响应
func Success[T any](data T) BaseResponse[T] {
	return NewBaseResponse(0, data, "ok")
}

// Error 失败响应（错误码）
func Error(errCode exception.ErrorCode) BaseResponse[any] {
	return NewBaseResponseFromErrorCode(errCode.Code(), errCode.Message())
}

// ErrorWithCode 失败响应（code + message）
func ErrorWithCode(code int, message string) BaseResponse[any] {
	return NewBaseResponseFromErrorCode(code, message)
}

// ErrorWithMessage 失败响应（错误码 + 自定义 message）
func ErrorWithMessage(errCode exception.ErrorCode, message string) BaseResponse[any] {
	return NewBaseResponseFromErrorCode(errCode.Code(), message)
}
