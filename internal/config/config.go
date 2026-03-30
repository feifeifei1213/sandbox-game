package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 汇总首版 Go 服务的基础配置。
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type ServerConfig struct {
	Host                string `mapstructure:"host"`
	Port                int    `mapstructure:"port"`
	Mode                string `mapstructure:"mode"`
	ReadTimeoutSeconds  int    `mapstructure:"readTimeoutSeconds"`
	WriteTimeoutSeconds int    `mapstructure:"writeTimeoutSeconds"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type MySQLConfig struct {
	DSN          string `mapstructure:"dsn"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
}

type AuthConfig struct {
	Mode               string `mapstructure:"mode"`
	TokenSecret        string `mapstructure:"tokenSecret"`
	TokenExpireSeconds int64  `mapstructure:"tokenExpireSeconds"`
}

// Address 返回 HTTP 服务监听地址。
func (c ServerConfig) Address() string {
	host := strings.TrimSpace(c.Host)
	if host == "" {
		host = "0.0.0.0"
	}
	if c.Port <= 0 {
		return fmt.Sprintf("%s:%d", host, 8080)
	}
	return fmt.Sprintf("%s:%d", host, c.Port)
}

// Load 从 YAML 配置文件读取最小运行配置。
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.SetDefault("app.name", "sandbox-game")
	v.SetDefault("app.version", "0.1.0")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.readTimeoutSeconds", 10)
	v.SetDefault("server.writeTimeoutSeconds", 10)
	v.SetDefault("log.level", "debug")
	v.SetDefault("mysql.dsn", "root:root@tcp(127.0.0.1:3306)/sandbox_game?charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("mysql.maxOpenConns", 10)
	v.SetDefault("mysql.maxIdleConns", 5)
	v.SetDefault("auth.mode", "local")
	v.SetDefault("auth.tokenSecret", "sandbox-game-local-secret")
	v.SetDefault("auth.tokenExpireSeconds", 28800)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
