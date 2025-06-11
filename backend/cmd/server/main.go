package main

import (
	"fmt"
	"log"
	"net/http"

	"temp-mailbox-service/internal/api"
	"temp-mailbox-service/internal/application"
	"temp-mailbox-service/internal/infrastructure/auth"
	"temp-mailbox-service/internal/infrastructure/config"
	"temp-mailbox-service/internal/infrastructure/database"
	"temp-mailbox-service/internal/infrastructure/middleware"
	"temp-mailbox-service/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
)

// 版本信息（构建时注入）
var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	fmt.Printf("临时邮箱系统 v%s\n", Version)
	fmt.Printf("构建时间: %s\n", BuildTime)
	fmt.Printf("Go版本: %s\n", GoVersion)

	// 加载配置
	cfg, err := config.Load(".env.dev")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	defer database.CloseDatabase()

	// 执行数据库迁移
	if err := database.Migrate(); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 创建服务实例
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)

	userRepo := persistence.NewUserRepository()
	userService := application.NewUserService(userRepo, jwtService)

	// 创建API处理器
	userHandler := api.NewUserHandler(userService)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin路由器
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// 健康检查端点（无需认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "temp-mailbox-service",
			"version":   Version,
			"buildTime": BuildTime,
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 公开的认证路由
		userHandler.RegisterRoutes(api)

		// 需要认证的用户路由
		userAuth := api.Group("/user")
		userAuth.Use(middleware.JWTAuth(jwtService))
		{
			userAuth.GET("/profile", userHandler.GetProfile)
			userAuth.PUT("/profile", userHandler.UpdateProfile)
			userAuth.POST("/change-password", userHandler.ChangePassword)
		}

		// 测试端点
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// 启动服务器
	addr := cfg.GetServerAddress()
	fmt.Printf("服务器启动在地址 %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
} 