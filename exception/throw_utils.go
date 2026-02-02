package exception

// ThrowIf 条件成立则抛出异常（返回 error）
func ThrowIf(condition bool, err error) error {
	if condition {
		return err
	}
	return nil
}

// ThrowIfWithCode 条件成立则返回错误码异常
func ThrowIfWithCode(condition bool, errCode ErrorCode) error {
	if condition {
		return NewBusinessErrorFromCode(errCode)
	}
	return nil
}

// ThrowIfWithMessage 条件成立则返回错误码异常（自定义消息）
func ThrowIfWithMessage(condition bool, errCode ErrorCode, message string) error {
	if condition {
		return NewBusinessErrorWithMessage(errCode, message)
	}
	return nil
}

// PanicIf 条件成立则 panic（用于无法处理的严重错误）
func PanicIf(condition bool, err error) {
	if condition {
		panic(err)
	}
}

// PanicIfWithCode 条件成立则 panic 业务错误
func PanicIfWithCode(condition bool, errCode ErrorCode) {
	if condition {
		panic(NewBusinessErrorFromCode(errCode))
	}
}
