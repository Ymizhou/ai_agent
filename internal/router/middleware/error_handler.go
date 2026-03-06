package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"aicode/internal/common"
	"aicode/internal/exception"
	applog "aicode/log"

	"github.com/gin-gonic/gin"
)

// GlobalErrorHandler 全局异常处理中间件
// 捕获 panic 时会附加完整调用堆栈和 traceId 到日志中
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				traceId := getTraceId(c)
				stack := applog.WithStack(err)

				var response common.BaseResponse[any]

				if bizErr, ok := err.(*exception.BusinessError); ok {
					applog.WithTraceId(c.Request.Context(), traceId).
						WithField("stack", stack).
						Errorf("panic recovered [BusinessError]: [%d] %s", bizErr.Code(), bizErr.Message())
					response = common.ErrorWithCode(bizErr.Code(), bizErr.Message())
				} else {
					applog.WithTraceId(c.Request.Context(), traceId).
						WithField("stack", stack).
						Errorf("panic recovered [SystemError]: %v", err)
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
	traceId := getTraceId(c)
	var response common.BaseResponse[any]

	var bizErr *exception.BusinessError
	if errors.As(err, &bizErr) {
		applog.WithTraceId(c.Request.Context(), traceId).
			Errorf("BusinessError: [%d] %s", bizErr.Code(), bizErr.Message())
		response = common.ErrorWithCode(bizErr.Code(), bizErr.Message())
	} else {
		applog.WithTraceId(c.Request.Context(), traceId).
			Errorf("SystemError: %v", err)
		response = common.ErrorWithCode(exception.SystemError.Code(), fmt.Sprintf("系统错误: %v", err))
	}

	c.JSON(http.StatusOK, response)
	c.Abort()
}

// getTraceId 从 gin.Context 中安全获取 traceId
func getTraceId(c *gin.Context) string {
	val, exists := c.Get(applog.TraceIdKey)
	if !exists {
		return ""
	}
	traceId, _ := val.(string)
	return traceId
}
