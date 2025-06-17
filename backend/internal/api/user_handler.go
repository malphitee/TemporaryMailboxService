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
			"code": 1001,
			"message": "请求参数格式错误",
			"data": nil,
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1002,
			"message": "请求参数验证失败",
			"data": nil,
		})
		return
	}
	
	// 调用用户服务注册
	loginResp, err := h.userService.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1003,
			"message": err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "用户注册成功",
		"data": loginResp,
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1101,
			"message": "请求参数格式错误",
			"data": nil,
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1102,
			"message": "请求参数验证失败",
			"data": nil,
		})
		return
	}
	
	// 调用用户服务登录
	loginResp, err := h.userService.LoginUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1103,
			"message": err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "登录成功",
		"data": loginResp,
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 1201,
			"message": "获取用户信息失败",
			"data": nil,
		})
		return
	}
	
	// 调用用户服务获取资料
	userResp, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1202,
			"message": err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "获取用户资料成功",
		"data": userResp,
	})
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 1301,
			"message": "获取用户信息失败",
			"data": nil,
		})
		return
	}
	
	var req user.UpdateUserRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1302,
			"message": "请求参数格式错误",
			"data": nil,
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1303,
			"message": "请求参数验证失败",
			"data": nil,
		})
		return
	}
	
	// 调用用户服务更新资料
	userResp, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1304,
			"message": err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "更新用户资料成功",
		"data": userResp,
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 1401,
			"message": "获取用户信息失败",
			"data": nil,
		})
		return
	}
	
	var req user.ChangePasswordRequest
	
	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1402,
			"message": "请求参数格式错误",
			"data": nil,
		})
		return
	}
	
	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1403,
			"message": "请求参数验证失败",
			"data": nil,
		})
		return
	}
	
	// 调用用户服务修改密码
	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1404,
			"message": err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "密码修改成功",
		"data": nil,
	})
} 