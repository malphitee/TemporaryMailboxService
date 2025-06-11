package api

import (
	"net/http"

	"temp-mailbox-service/internal/application"
	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService application.UserService
	validator   *validator.Validate
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数格式错误",
			"code":  "INVALID_REQUEST_FORMAT",
			"details": err.Error(),
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败",
			"code":  "VALIDATION_FAILED",
			"details": err.Error(),
		})
		return
	}
	
	// 调用用户服务注册
	userResp, err := h.userService.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "REGISTRATION_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"data":    userResp,
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数格式错误",
			"code":  "INVALID_REQUEST_FORMAT",
			"details": err.Error(),
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败",
			"code":  "VALIDATION_FAILED",
			"details": err.Error(),
		})
		return
	}
	
	// 调用用户服务登录
	loginResp, err := h.userService.LoginUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
			"code":  "LOGIN_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    loginResp,
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户信息失败",
			"code":  "UNAUTHORIZED",
		})
		return
	}
	
	// 调用用户服务获取资料
	userResp, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
			"code":  "USER_NOT_FOUND",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户资料成功",
		"data":    userResp,
	})
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户信息失败",
			"code":  "UNAUTHORIZED",
		})
		return
	}
	
	var req user.UpdateUserRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数格式错误",
			"code":  "INVALID_REQUEST_FORMAT",
			"details": err.Error(),
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败",
			"code":  "VALIDATION_FAILED",
			"details": err.Error(),
		})
		return
	}
	
	// 调用用户服务更新资料
	userResp, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "UPDATE_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "更新用户资料成功",
		"data":    userResp,
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户信息失败",
			"code":  "UNAUTHORIZED",
		})
		return
	}
	
	var req user.ChangePasswordRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数格式错误",
			"code":  "INVALID_REQUEST_FORMAT",
			"details": err.Error(),
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败",
			"code":  "VALIDATION_FAILED",
			"details": err.Error(),
		})
		return
	}
	
	// 调用用户服务修改密码
	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "CHANGE_PASSWORD_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功",
	})
} 