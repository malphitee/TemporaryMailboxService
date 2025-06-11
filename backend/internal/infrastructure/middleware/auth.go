package middleware

import (
	"net/http"
	"strconv"
	"time"

	"temp-mailbox-service/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "缺少授权头",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 提取Bearer token
		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "无效的授权格式",
				"message": "Authorization header must be in format 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// 验证访问令牌
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "令牌验证失败",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息设置到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.Next()
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			// 可选认证失败时不阻止请求，但记录日志
			c.Header("X-Auth-Warning", "Invalid token provided")
			c.Next()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

// GetCurrentUsername 从上下文中获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// GetCurrentUserEmail 从上下文中获取当前用户邮箱
func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	userEmail, ok := email.(string)
	return userEmail, ok
}

// GetJWTClaims 从上下文中获取JWT声明
func GetJWTClaims(c *gin.Context) (*auth.JWTClaims, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}
	jwtClaims, ok := claims.(*auth.JWTClaims)
	return jwtClaims, ok
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以根据实际需求检查用户角色
		// 目前简化处理，只检查是否已认证
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "需要认证",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// TODO: 在用户模型中添加角色字段后，可以在这里检查管理员权限
		// 目前所有已认证用户都被视为管理员

		c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

// RateLimiterMiddleware 简单的速率限制中间件
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现基于IP或用户的速率限制
		// 目前暂时跳过
		c.Next()
	}
}

// JWTAuth JWT认证中间件
func JWTAuth(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
				"code":  "MISSING_TOKEN",
			})
			c.Abort()
			return
		}
		
		// 提取Bearer令牌
		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "令牌格式错误",
				"code":  "INVALID_TOKEN_FORMAT",
			})
			c.Abort()
			return
		}
		
		// 验证令牌
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "令牌无效或已过期",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}
		
		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		
		// 继续处理请求
		c.Next()
	}
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePublic}
	}
	
	if id, ok := userID.(uint); ok {
		return id, nil
	}
	
	return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePublic}
}

// GetUsernameFromContext 从上下文中获取用户名
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	
	if name, ok := username.(string); ok {
		return name, true
	}
	
	return "", false
}

// GetEmailFromContext 从上下文中获取邮箱
func GetEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	
	if addr, ok := email.(string); ok {
		return addr, true
	}
	
	return "", false
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func OptionalAuth(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有认证头，继续处理但不设置用户信息
			c.Next()
			return
		}
		
		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			// 令牌格式错误，继续处理但不设置用户信息
			c.Next()
			return
		}
		
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			// 令牌无效，继续处理但不设置用户信息
			c.Next()
			return
		}
		
		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		
		c.Next()
	}
}

// RequireRole 角色验证中间件（为将来的角色系统预留）
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以实现角色验证逻辑
		// 目前只是一个占位符
		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		}
		
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RateLimit 限流中间件（简单版本）
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以实现限流逻辑
		// 目前只是一个占位符
		c.Next()
	}
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 生成简单的请求ID（使用时间戳）
			requestID = strconv.FormatInt(time.Now().Unix(), 36)
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
} 