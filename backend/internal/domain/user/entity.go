package user

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 基本信息
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
	Email    string `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email,max=100"`
	Password string `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
	
	// 用户状态
	IsActive bool `json:"is_active" gorm:"default:true"`
	
	// 用户资料
	FirstName string `json:"first_name" gorm:"size:50" validate:"max=50"`
	LastName  string `json:"last_name" gorm:"size:50" validate:"max=50"`
	Avatar    string `json:"avatar" gorm:"size:255"`
	
	// 认证相关
	LastLoginAt    *time.Time `json:"last_login_at"`
	PasswordResetToken string `json:"-" gorm:"size:255"`
	PasswordResetExpiry *time.Time `json:"-"`
	
	// 用户设置
	TimeZone string `json:"timezone" gorm:"size:50;default:'UTC'"`
	Language string `json:"language" gorm:"size:10;default:'zh-CN'"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// GetFullName 获取用户全名
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

// IsPasswordResetValid 检查密码重置令牌是否有效
func (u *User) IsPasswordResetValid() bool {
	if u.PasswordResetToken == "" || u.PasswordResetExpiry == nil {
		return false
	}
	return time.Now().Before(*u.PasswordResetExpiry)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"max=50"`
	LastName  string `json:"last_name" validate:"max=50"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"max=50"`
	LastName  string `json:"last_name" validate:"max=50"`
	Avatar    string `json:"avatar" validate:"max=255"`
	TimeZone  string `json:"timezone" validate:"max=50"`
	Language  string `json:"language" validate:"max=10"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	IsActive    bool       `json:"is_active"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Avatar      string     `json:"avatar"`
	LastLoginAt *time.Time `json:"last_login_at"`
	TimeZone    string     `json:"timezone"`
	Language    string     `json:"language"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Username:    u.Username,
		Email:       u.Email,
		IsActive:    u.IsActive,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Avatar:      u.Avatar,
		LastLoginAt: u.LastLoginAt,
		TimeZone:    u.TimeZone,
		Language:    u.Language,
	}
} 