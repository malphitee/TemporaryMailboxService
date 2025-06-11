package main

import (
	"fmt"
	"log"
	"net/http"

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

	// 创建Gin路由器
	r := gin.Default()

	// 健康检查端点
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
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// 启动服务器
	port := ":8080"
	fmt.Printf("服务器启动在端口 %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
} 