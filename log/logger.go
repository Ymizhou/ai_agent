package log

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

const (
	// TraceIdHeader 前端携带 traceId 的请求头名称
	TraceIdHeader = "X-Trace-Id"
	// TraceIdKey gin.Context 中存储 traceId 的键名
	TraceIdKey = "trace_id"
)

// logger 全局 logrus 实例
var logger = logrus.New()

// Init 初始化日志模块，应在应用启动时调用一次
// level 对应配置文件中的 server.log_level，例如 "debug"、"info"、"warn"、"error"
func Init(level string) {
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
		},
	})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Warnf("无效的日志级别 %q，使用默认级别 info", level)
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// 同步到全局 logrus，兼容项目中已有的 logrus.Xxx() 调用
	logrus.SetOutput(logger.Out)
	logrus.SetFormatter(logger.Formatter)
	logrus.SetLevel(logger.Level)
}

// WithTraceId 返回携带 trace_id 字段的 Entry，用于链路追踪日志
func WithTraceId(ctx context.Context, traceId string) *logrus.Entry {
	if traceId != "" {
		return logger.WithField("trace_id", traceId)
	}
	return logger.WithField("trace_id", GetTraceId(ctx))
}

func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(TraceIdKey).(string)
	if !ok {
		return ""
	}
	return traceId
}

// WithStack 将 recover() 捕获的 panic 值与当前调用堆栈合并为一个字段串，
// 供 error_handler 中间件附加到日志上。
func WithStack(recovered interface{}) string {
	return string(debug.Stack())
}

// GetLogger 返回全局 Logger 实例，供其他包直接使用
func GetLogger() *logrus.Logger {
	return logger
}
