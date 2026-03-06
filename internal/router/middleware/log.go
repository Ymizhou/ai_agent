package middleware

import (
	"aicode/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func FormatLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()
		logrus.SetFormatter(&logrus.TextFormatter{})
		logLevel, err := logrus.ParseLevel(cfg.Server.LogLevel)
		if err != nil {
			logrus.Panicf("解析日志级别失败: %s, 使用默认级别: %s", err.Error(), logrus.DebugLevel)
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logLevel)
		}
		c.Next()
	}
}
