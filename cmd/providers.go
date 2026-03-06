package cmd

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"aicode/ai/chatmodel"
	"aicode/config"
)

// ProvideConfig 提供配置
func MustProvideConfig() *config.Config {
	// 加载配置文件
	return config.LoadConfig("config.yml")
}

// ProvideDB 提供数据库实例
func MustProvideDB(cfg *config.Config) *gorm.DB {
	config.InitDatabase()
	return config.GetDB()
}

// ProvideChatModel 提供聊天模型实例
func MustProvideChatModel(cfg *config.Config) (
	map[string]chatmodel.ChatModelFactory,
	error,
) {
	chatModelRegistry, err := chatmodel.InitChatModel(cfg)
	if err != nil {
		logrus.Panicf("初始化聊天模型失败: %v", err)
	}
	return chatModelRegistry, nil
}
