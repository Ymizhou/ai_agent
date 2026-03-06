package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	applog "aicode/log"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装 gin.ResponseWriter，拦截响应体内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// AccessLogMiddleware 请求访问日志中间件
// 记录每个请求的：method、path、query、client_ip、request_body、
// 以及响应的：status、latency_ms、response_body（SSE 流式接口跳过响应体）
func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读取请求体（读完后需要重新填回，否则后续 Handler 无法再次读取）
		var reqBodyStr string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				reqBodyStr = string(bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 判断是否为 SSE 流式接口，流式接口不捕获响应体
		isSSE := strings.Contains(c.Request.URL.Path, "/stream")

		// 包装 ResponseWriter 以捕获响应体（非流式）
		var rw *responseWriter
		if !isSSE {
			rw = &responseWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = rw
		}

		// 取 traceId
		traceId, _ := c.Get(applog.TraceIdKey)
		traceIdStr, _ := traceId.(string)

		// 记录请求进入日志
		applog.WithTraceId(c.Request.Context(), traceIdStr).
			WithField("phase", "request").
			WithField("method", c.Request.Method).
			WithField("path", c.Request.URL.Path).
			WithField("query", c.Request.URL.RawQuery).
			WithField("client_ip", c.ClientIP()).
			WithField("request_body", reqBodyStr).
			Info("incoming request")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// 构建响应日志 Entry
		entry := applog.WithTraceId(c.Request.Context(), traceIdStr).
			WithField("phase", "response").
			WithField("method", c.Request.Method).
			WithField("path", c.Request.URL.Path).
			WithField("status", status).
			WithField("latency_ms", latency.Milliseconds())

		if isSSE {
			entry = entry.WithField("response_body", "[SSE stream]")
		} else if rw != nil {
			entry = entry.WithField("response_body", rw.body.String())
		}

		// 有错误时附加错误信息
		if len(c.Errors) > 0 {
			entry = entry.WithField("errors", c.Errors.Errors())
			entry.Error("request completed with errors")
		} else {
			entry.Info("request completed")
		}
	}
}
