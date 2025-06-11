package user

import (
	"context"
	"time"
)

// Repository 用户仓储接口
type Repository interface {
	// 基础CRUD操作
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	List(ctx context.Context, offset, limit int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context, id uint) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// 认证相关
	UpdatePassword(ctx context.Context, id uint, hashedPassword string) error
	UpdateLastLogin(ctx context.Context, id uint) error
	SetPasswordResetToken(ctx context.Context, id uint, token string, expiry *time.Time) error
	ClearPasswordResetToken(ctx context.Context, id uint) error
	
	// 状态管理
	Activate(ctx context.Context, id uint) error
	Deactivate(ctx context.Context, id uint) error
} 