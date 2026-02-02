package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	cfg := GetConfig()
	dsn := cfg.Database.GetDSN()

	// 配置 GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		logrus.Panicf("连接数据库失败: %v", err)
	}

	// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Panicf("获取数据库连接失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数
	// sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		logrus.Panicf("数据库连接测试失败: %v", err)
	}

	DB = db
	logrus.Info("数据库连接成功")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
