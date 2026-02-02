package exception

// ErrorCode 错误码接口
type ErrorCode interface {
	Code() int
	Message() string
	Error() string // 实现 error 接口
}

// BaseErrorCode 基础错误码实现
type BaseErrorCode struct {
	code    int
	message string
}

func (e BaseErrorCode) Code() int {
	return e.code
}

func (e BaseErrorCode) Message() string {
	return e.message
}

func (e BaseErrorCode) Error() string {
	return e.message
}

// 定义常用错误码
var (
	SUCCESS           = BaseErrorCode{0, "ok"}
	ParamsError       = BaseErrorCode{40000, "请求参数错误"}
	NotLoginError     = BaseErrorCode{40100, "未登录"}
	NoAuthError       = BaseErrorCode{40101, "无权限"}
	TooManyRequest    = BaseErrorCode{42900, "请求过于频繁"}
	NotFoundError     = BaseErrorCode{40400, "请求数据不存在"}
	ForbiddenError    = BaseErrorCode{40300, "禁止访问"}
	SystemError       = BaseErrorCode{50000, "系统内部异常"}
	OperationError    = BaseErrorCode{50001, "操作失败"}
)
