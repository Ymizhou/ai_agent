package main

import (
	"aicode/config"
	"database/sql"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("config", "config.yml", "配置文件路径")
	migrationsPath := flag.String("migrations", "file://migrations", "迁移文件目录路径")
	flag.Parse()

	// 加载配置文件
	config.LoadConfig(*configPath)

	cfg := config.GetConfig()
	logrus.Infof("开始数据库迁移，数据库: %s", cfg.Database.DBName)

	// 连接数据库
	db, err := sql.Open("mysql", cfg.Database.GetDSN())
	if err != nil {
		logrus.Panicf("连接数据库失败: %s", err.Error())
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		logrus.Panicf("数据库连接测试失败: %s", err.Error())
	}
	logrus.Info("数据库连接成功")

	// 创建 mysql driver 实例
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		logrus.Panicf("创建数据库驱动失败: %s", err.Error())
	}

	// 创建 migrate 实例
	m, err := migrate.NewWithDatabaseInstance(
		*migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		logrus.Panicf("创建迁移实例失败: %s", err.Error())
	}

	// 获取当前版本
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logrus.Panicf("获取当前版本失败: %s", err.Error())
	} else if err == migrate.ErrNilVersion {
		logrus.Info("当前版本: 无（首次迁移）")
		logrus.Info("即将创建版本管理表 schema_migrations 并执行迁移...")
	} else {
		logrus.Infof("当前版本: %d, 脏数据状态: %v", currentVersion, dirty)

		// 如果存在脏数据状态，需要手动修复
		if dirty {
			logrus.Panicf("警告: 数据库处于脏数据状态，可能是上次迁移失败导致，请检查数据库状态，必要时使用 -action=force 强制设置版本")
		}
	}

	// 执行迁移
	logrus.Info("开始执行数据库迁移...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logrus.Info("数据库已是最新版本，无需迁移")
		} else {
			logrus.Panicf("迁移失败: %v", err)
		}
	} else {
		logrus.Info("迁移成功完成！")
	}

	// 显示最终版本
	finalVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logrus.Panicf("获取最终版本失败: %v", err)
	} else if err == migrate.ErrNilVersion {
		logrus.Info("最终版本: 无")
	} else {
		logrus.Infof("最终版本: %d, 脏数据状态: %v", finalVersion, dirty)
	}
}
