package middleware

import (
	"aicode/common"
	"aicode/exception"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GlobalErrorHandler 全局异常处理中间件
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 处理 panic
				logrus.Errorf("Panic recovered: %v", err)

				var response common.BaseResponse[any]

				// 判断是否是业务异常
				if bizErr, ok := err.(*exception.BusinessError); ok {
					logrus.Errorf("BusinessError: [%d] %s", bizErr.Code(), bizErr.Message())
					response = common.ErrorWithCode(bizErr.Code(), bizErr.Message())
				} else {
					// 其他类型的 panic
					logrus.Errorf("SystemError: %v", err)
					response = common.Error(exception.SystemError)
				}

				c.JSON(http.StatusOK, response)
				c.Abort()
			}
		}()

		c.Next()

		// 检查是否有错误被设置到上下文中
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

// handleError 处理错误并返回响应
func handleError(c *gin.Context, err error) {
	var response common.BaseResponse[any]

	// 判断错误类型
	var bizErr *exception.BusinessError
	if errors.As(err, &bizErr) {
		// 业务异常
		logrus.Errorf("BusinessError: [%d] %s", bizErr.Code(), bizErr.Message())
		response = common.ErrorWithCode(bizErr.Code(), bizErr.Message())
	} else {
		// 其他错误
		logrus.Errorf("Error: %v", err)
		response = common.ErrorWithCode(exception.SystemError.Code(), fmt.Sprintf("系统错误: %v", err))
	}

	c.JSON(http.StatusOK, response)
	c.Abort()
}
