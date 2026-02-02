package cmd

import (
	"aicode/config"

	"gorm.io/gorm"
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
