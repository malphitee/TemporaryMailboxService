package config

import (
	"os"
	"testing"
)

func TestLoad_WithDefaults(t *testing.T) {
	// 测试使用默认配置加载
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("加载默认配置失败: %v", err)
	}

	// 验证默认值
	if cfg.Server.Host != "localhost" {
		t.Errorf("期望主机为 'localhost'，得到 '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("期望端口为 8080，得到 %d", cfg.Server.Port)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("期望数据库驱动为 'sqlite'，得到 '%s'", cfg.Database.Driver)
	}
	if cfg.JWT.AccessTokenTTL != 60 {
		t.Errorf("期望访问令牌TTL为 60，得到 %d", cfg.JWT.AccessTokenTTL)
	}
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// 设置环境变量
	os.Setenv("TEMP_MAILBOX_SERVER_PORT", "9090")
	os.Setenv("TEMP_MAILBOX_DATABASE_DRIVER", "postgres")
	os.Setenv("TEMP_MAILBOX_JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("TEMP_MAILBOX_SERVER_PORT")
		os.Unsetenv("TEMP_MAILBOX_DATABASE_DRIVER")
		os.Unsetenv("TEMP_MAILBOX_JWT_SECRET")
	}()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("加载环境变量配置失败: %v", err)
	}

	// 验证环境变量生效
	if cfg.Server.Port != 9090 {
		t.Errorf("期望端口为 9090，得到 %d", cfg.Server.Port)
	}
	if cfg.Database.Driver != "postgres" {
		t.Errorf("期望数据库驱动为 'postgres'，得到 '%s'", cfg.Database.Driver)
	}
	if cfg.JWT.Secret != "test-secret-key" {
		t.Errorf("期望JWT密钥为 'test-secret-key'，得到 '%s'", cfg.JWT.Secret)
	}
}

func TestValidateConfig_InvalidPort(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: -1,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "./test.db",
		},
		JWT: JWTConfig{
			Secret:          "test-secret",
			AccessTokenTTL:  60,
			RefreshTokenTTL: 10080,
		},
		Log: LogConfig{
			Level: "info",
		},
	}

	err := validateConfig(cfg)
	if err == nil {
		t.Error("无效端口应该导致验证失败")
	}
}

func TestValidateConfig_InvalidDriver(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver: "mysql", // 不支持的驱动
			DSN:    "./test.db",
		},
		JWT: JWTConfig{
			Secret:          "test-secret",
			AccessTokenTTL:  60,
			RefreshTokenTTL: 10080,
		},
		Log: LogConfig{
			Level: "info",
		},
	}

	err := validateConfig(cfg)
	if err == nil {
		t.Error("不支持的数据库驱动应该导致验证失败")
	}
}

func TestValidateConfig_ProductionWithDefaultSecret(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "release", // 生产模式
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "./test.db",
		},
		JWT: JWTConfig{
			Secret:          "your-secret-key-change-in-production", // 默认密钥
			AccessTokenTTL:  60,
			RefreshTokenTTL: 10080,
		},
		Log: LogConfig{
			Level: "info",
		},
	}

	err := validateConfig(cfg)
	if err == nil {
		t.Error("生产环境使用默认JWT密钥应该导致验证失败")
	}
}

func TestValidateConfig_Valid(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "./test.db",
		},
		JWT: JWTConfig{
			Secret:          "secure-secret-key",
			AccessTokenTTL:  60,
			RefreshTokenTTL: 10080,
		},
		Log: LogConfig{
			Level: "info",
		},
	}

	err := validateConfig(cfg)
	if err != nil {
		t.Errorf("有效配置验证失败: %v", err)
	}
}

func TestIsProduction(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Mode: "release"},
	}
	if !cfg.IsProduction() {
		t.Error("release模式应该被识别为生产环境")
	}

	cfg.Server.Mode = "debug"
	if cfg.IsProduction() {
		t.Error("debug模式不应该被识别为生产环境")
	}
}

func TestGetServerAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	addr := cfg.GetServerAddress()
	expected := "localhost:8080"
	if addr != expected {
		t.Errorf("期望服务器地址为 '%s'，得到 '%s'", expected, addr)
	}
} 