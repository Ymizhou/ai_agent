package main

import (
	"aicode/cmd"
	"aicode/config"
	applog "aicode/log"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置（日志初始化依赖配置中的 log_level，需先加载）
	cfg := config.GetConfig()
	if cfg == nil {
		// 配置尚未加载时使用默认 info 级别，待配置加载后 Init 会覆盖
		applog.Init("info")
	} else {
		applog.Init(cfg.Server.LogLevel)
	}

	// 使用Wire初始化应用程序
	app, err := cmd.InitializeApp()
	if err != nil {
		logrus.Panicf("初始化应用失败: %v", err)
	}

	// 配置加载完成后重新初始化日志（InitializeApp 内部会调用 config.LoadConfig）
	cfg = config.GetConfig()
	applog.Init(cfg.Server.LogLevel)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("服务器启动在端口: %d, 根路径: %s", cfg.Server.Port, cfg.Server.RootPath)
	if err := app.Router.Run(addr); err != nil {
		logrus.Panicf("服务器启动失败: %v", err)
	}
}
