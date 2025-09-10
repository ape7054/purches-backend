package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	App      AppConfig      `mapstructure:"app"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
	FilePath string `mapstructure:"file_path"` // SQLite文件路径
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
	DefaultUser string `mapstructure:"default_user_id"`
}

var globalConfig *Config

// LoadConfig 加载配置
func LoadConfig(configPath ...string) (*Config, error) {
	config := &Config{}

	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 添加配置文件搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// 如果指定了配置路径
	if len(configPath) > 0 {
		viper.SetConfigFile(configPath[0])
	}

	// 设置环境变量前缀
	viper.SetEnvPrefix("PURCHES")
	viper.AutomaticEnv()

	// 环境变量键名替换
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认配置
		fmt.Printf("配置文件未找到，使用默认配置: %v\n", err)
	} else {
		fmt.Printf("使用配置文件: %s\n", viper.ConfigFileUsed())
	}

	// 将配置解析到结构体
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 全局配置
	globalConfig = config

	return config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", 60)
	viper.SetDefault("server.write_timeout", 60)

	// 数据库配置
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.file_path", "./purches.db")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "purches")
	viper.SetDefault("database.ssl_mode", "disable")

	// 应用配置
	viper.SetDefault("app.name", "采购订单系统")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", getEnv("ENVIRONMENT", "development"))
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.default_user_id", "user_1")
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		// 如果配置未初始化，尝试加载默认配置
		config, err := LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("配置加载失败: %v", err))
		}
		return config
	}
	return globalConfig
}

// getEnv 获取环境变量，带默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	return GetConfig().App.Environment == "development"
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	return GetConfig().App.Environment == "production"
}

// GetDatabaseDSN 获取数据库连接字符串
func GetDatabaseDSN() string {
	cfg := GetConfig().Database

	switch cfg.Type {
	case "sqlite":
		return cfg.FilePath
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	default:
		return cfg.FilePath // 默认使用SQLite
	}
}
