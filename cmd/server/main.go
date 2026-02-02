package main

import (
	"aicode/cmd"
	"aicode/config"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	// 使用Wire初始化应用程序
	app, err := cmd.InitializeApp()
	if err != nil {
		logrus.Panicf("初始化应用失败: %v", err)
	}
	// 启动服务器
	cfg := config.GetConfig()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("服务器启动在端口: %d, 上下文路径: %s", cfg.Server.Port, cfg.Server.ContextPath)
	if err := app.Router.Run(addr); err != nil {
		logrus.Panicf("服务器启动失败: %v", err)
	}
}
