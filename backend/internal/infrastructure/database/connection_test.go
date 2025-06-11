package database

import (
	"os"
	"testing"

	"temp-mailbox-service/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDatabase(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.DatabaseConfig
		wantErr bool
	}{
		{
			name: "SQLite数据库配置正确",
			config: &config.DatabaseConfig{
				Driver:       "sqlite",
				DSN:          "./test.db",
				MaxOpenConns: 25,
				MaxIdleConns: 10,
				MaxLifetime:  30,
			},
			wantErr: false,
		},
		{
			name: "内存SQLite数据库",
			config: &config.DatabaseConfig{
				Driver:       "sqlite",
				DSN:          ":memory:",
				MaxOpenConns: 25,
				MaxIdleConns: 10,
				MaxLifetime:  30,
			},
			wantErr: false,
		},
		{
			name: "不支持的数据库驱动",
			config: &config.DatabaseConfig{
				Driver:       "mysql",
				DSN:          "test",
				MaxOpenConns: 25,
				MaxIdleConns: 10,
				MaxLifetime:  30,
			},
			wantErr: true,
		},
		{
			name: "PostgreSQL配置（模拟）",
			config: &config.DatabaseConfig{
				Driver:       "postgres",
				DSN:          "postgres://user:pass@localhost/test?sslmode=disable",
				MaxOpenConns: 25,
				MaxIdleConns: 10,
				MaxLifetime:  30,
			},
			wantErr: true, // 在测试环境中PostgreSQL不可用
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理之前的全局DB实例
			DB = nil
			
			err := InitDatabase(tt.config)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, DB)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, DB)
				
				// 测试数据库连接
				sqlDB, err := DB.DB()
				assert.NoError(t, err)
				assert.NoError(t, sqlDB.Ping())
				
				// 验证连接池配置
				assert.Equal(t, tt.config.MaxOpenConns, sqlDB.Stats().MaxOpenConnections)
				
				// 清理
				CloseDatabase()
			}
		})
	}
}

func TestGetDB(t *testing.T) {
	// 测试未初始化时的情况
	DB = nil
	db := GetDB()
	assert.Nil(t, db)
	
	// 初始化数据库
	cfg := &config.DatabaseConfig{
		Driver:       "sqlite",
		DSN:          ":memory:",
		MaxOpenConns: 25,
		MaxIdleConns: 10,
		MaxLifetime:  30,
	}
	
	err := InitDatabase(cfg)
	require.NoError(t, err)
	
	// 测试获取数据库实例
	db = GetDB()
	assert.NotNil(t, db)
	assert.Equal(t, DB, db)
	
	// 清理
	CloseDatabase()
}

func TestCloseDatabase(t *testing.T) {
	tests := []struct {
		name      string
		setupDB   bool
		expectErr bool
	}{
		{
			name:      "关闭已初始化的数据库",
			setupDB:   true,
			expectErr: false,
		},
		{
			name:      "关闭未初始化的数据库",
			setupDB:   false,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理状态
			DB = nil
			
			if tt.setupDB {
				cfg := &config.DatabaseConfig{
					Driver:       "sqlite",
					DSN:          ":memory:",
					MaxOpenConns: 25,
					MaxIdleConns: 10,
					MaxLifetime:  30,
				}
				err := InitDatabase(cfg)
				require.NoError(t, err)
			}
			
			err := CloseDatabase()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatabaseDriverSupport(t *testing.T) {
	supportedDrivers := []string{"sqlite", "postgres"}
	unsupportedDrivers := []string{"mysql", "oracle", "sqlserver", ""}

	// 测试支持的驱动
	for _, driver := range supportedDrivers {
		t.Run("支持的驱动_"+driver, func(t *testing.T) {
			cfg := &config.DatabaseConfig{
				Driver:       driver,
				DSN:          getTestDSN(driver),
				MaxOpenConns: 25,
				MaxIdleConns: 10,
				MaxLifetime:  30,
			}
			
			DB = nil
			err := InitDatabase(cfg)
			
			if driver == "sqlite" {
				// SQLite应该总是成功
				assert.NoError(t, err)
				assert.NotNil(t, DB)
				CloseDatabase()
			} else {
				// PostgreSQL在测试环境中可能失败，但应该支持驱动
				// 错误应该是连接错误，不是驱动不支持错误
				if err != nil {
					assert.NotContains(t, err.Error(), "不支持的数据库驱动")
				}
			}
		})
	}

	// 测试不支持的驱动
	for _, driver := range unsupportedDrivers {
		t.Run("不支持的驱动_"+driver, func(t *testing.T) {
			cfg := &config.DatabaseConfig{
				Driver:       driver,
				DSN:          "test",
				MaxOpenConns: 25,
				MaxIdleConns: 10,
				MaxLifetime:  30,
			}
			
			DB = nil
			err := InitDatabase(cfg)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "不支持的数据库驱动")
			assert.Nil(t, DB)
		})
	}
}

func TestConnectionPoolConfiguration(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Driver:       "sqlite",
		DSN:          ":memory:",
		MaxOpenConns: 50,
		MaxIdleConns: 20,
		MaxLifetime:  60,
	}
	
	DB = nil
	err := InitDatabase(cfg)
	require.NoError(t, err)
	
	sqlDB, err := DB.DB()
	require.NoError(t, err)
	
	// 验证连接池配置
	stats := sqlDB.Stats()
	assert.Equal(t, cfg.MaxOpenConns, stats.MaxOpenConnections)
	
	// 注意：Go的sql.DB不直接暴露MaxIdleConns和ConnMaxLifetime的getter
	// 但我们可以验证配置被应用了（通过没有错误来间接验证）
	assert.NoError(t, sqlDB.Ping())
	
	CloseDatabase()
}

// getTestDSN 根据驱动类型返回测试用的DSN
func getTestDSN(driver string) string {
	switch driver {
	case "sqlite":
		return ":memory:"
	case "postgres":
		return "postgres://test:test@localhost/test?sslmode=disable"
	default:
		return "test"
	}
}

// 清理测试文件
func TestMain(m *testing.M) {
	// 运行测试
	code := m.Run()
	
	// 清理测试数据库文件
	os.Remove("./test.db")
	
	os.Exit(code)
} 