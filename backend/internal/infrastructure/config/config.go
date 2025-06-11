package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`         // debug, release
	ReadTimeout  int    `mapstructure:"read_timeout"` // seconds
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`        // sqlite, postgres
	DSN          string `mapstructure:"dsn"`          // 数据源名称
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"` // minutes
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	AccessTokenTTL   int    `mapstructure:"access_token_ttl"`   // minutes
	RefreshTokenTTL  int    `mapstructure:"refresh_token_ttl"`  // minutes
	Issuer           string `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, text
	Output string `mapstructure:"output"` // stdout, stderr, file path
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	v := viper.New()
	
	// 设置配置文件类型
	v.SetConfigType("env")
	
	// 设置环境变量前缀
	v.SetEnvPrefix("TEMP_MAILBOX")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	
	// 设置默认值
	setDefaults(v)
	
	// 如果提供了配置文件路径，则加载配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}
	
	// 解析配置到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	
	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}
	
	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 服务器默认配置
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	
	// 数据库默认配置
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.dsn", "./dev.db")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_lifetime", 30)
	
	// JWT默认配置
	v.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	v.SetDefault("jwt.access_token_ttl", 60)  // 1小时
	v.SetDefault("jwt.refresh_token_ttl", 10080) // 7天
	v.SetDefault("jwt.issuer", "temp-mailbox-service")
	
	// 日志默认配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.output", "stdout")
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("无效的服务器端口: %d", config.Server.Port)
	}
	
	if config.Server.Mode != "debug" && config.Server.Mode != "release" {
		return fmt.Errorf("无效的服务器模式: %s", config.Server.Mode)
	}
	
	// 验证数据库配置
	if config.Database.Driver != "sqlite" && config.Database.Driver != "postgres" {
		return fmt.Errorf("不支持的数据库驱动: %s", config.Database.Driver)
	}
	
	if config.Database.DSN == "" {
		return fmt.Errorf("数据库DSN不能为空")
	}
	
	// 验证JWT配置
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}
	// 在生产环境中警告使用默认密钥
	if config.JWT.Secret == "your-secret-key-change-in-production" && config.Server.Mode == "release" {
		return fmt.Errorf("生产环境中不能使用默认JWT密钥，请设置安全的密钥")
	}
	
	if config.JWT.AccessTokenTTL <= 0 {
		return fmt.Errorf("访问令牌TTL必须大于0")
	}
	
	// 验证日志配置
	validLogLevels := []string{"debug", "info", "warn", "error"}
	validLevel := false
	for _, level := range validLogLevels {
		if config.Log.Level == level {
			validLevel = true
			break
		}
	}
	if !validLevel {
		return fmt.Errorf("无效的日志级别: %s", config.Log.Level)
	}
	
	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return c.DSN
}

// IsProduction 检查是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Server.Mode == "release"
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
} 