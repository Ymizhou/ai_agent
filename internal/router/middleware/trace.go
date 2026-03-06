package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	applog "aicode/log"

	"github.com/gin-gonic/gin"
)

// TraceMiddleware 链路追踪中间件
// 优先读取请求头 X-Trace-Id，若前端未携带则自动生成随机 traceId。
// traceId 会写入 gin.Context 供后续中间件/Handler 使用，同时回写到响应头。
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(applog.TraceIdHeader)
		if traceId == "" {
			traceId = generateTraceId()
		}
		// 存入 Context，供后续中间件读取
		c.Set(applog.TraceIdKey, traceId)
		c.Request.WithContext(context.WithValue(c.Request.Context(), applog.TraceIdKey, traceId))
		// 回写响应头，方便前端关联追踪
		c.Header(applog.TraceIdHeader, traceId)
		c.Next()
	}
}

// generateTraceId 生成随机 traceId（时间戳 + 随机数，16进制格式）
func generateTraceId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d%08x", time.Now().UnixMilli(), r.Int31())
}
