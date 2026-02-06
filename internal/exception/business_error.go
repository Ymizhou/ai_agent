package exception

import "fmt"

// BusinessError 业务异常，实现 error 接口
type BusinessError struct {
	code    int
	message string
}

// NewBusinessError 创建业务异常
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		code:    code,
		message: message,
	}
}

// NewBusinessErrorFromCode 从错误码创建业务异常
func NewBusinessErrorFromCode(errCode ErrorCode) *BusinessError {
	return &BusinessError{
		code:    errCode.Code(),
		message: errCode.Message(),
	}
}

// NewBusinessErrorWithMessage 从错误码创建业务异常（自定义消息）
func NewBusinessErrorWithMessage(errCode ErrorCode, message string) *BusinessError {
	return &BusinessError{
		code:    errCode.Code(),
		message: message,
	}
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%d] %s", e.code, e.message)
}

// Code 获取错误码
func (e *BusinessError) Code() int {
	return e.code
}

// Message 获取错误消息
func (e *BusinessError) Message() string {
	return e.message
}
