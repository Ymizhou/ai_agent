package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var globalConfig *Config

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int    `yaml:"port"`
	ContextPath string `yaml:"context_path"`
	LogLevel    string `yaml:"log_level"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `yaml:"name"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parseTime"`
	Loc       string `yaml:"loc"`
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		d.Username, d.Password, d.Host, d.Port, d.DBName, d.Charset, d.ParseTime, d.Loc)
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) *Config {
	logrus.Infof("加载配置文件: %s", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Panicf("读取配置文件失败: %s", err.Error())
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logrus.Panicf("解析配置文件失败: %s", err.Error())
	}

	globalConfig = &cfg
	logrus.Infof("加载配置文件成功: %s", configPath)
	return &cfg
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}
