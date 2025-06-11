package persistence

import (
	"context"
	"errors"
	"time"

	"temp-mailbox-service/internal/domain/user"

	"gorm.io/gorm"
)

// userRepository GORM用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建新的用户仓储实例
func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

// Create 创建新用户
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var users []*user.User
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

// Count 获取用户总数
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error
	return count, err
}

// Exists 检查用户是否存在
func (r *userRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否已存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername 检查用户名是否已存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// SetPasswordResetToken 设置密码重置令牌
func (r *userRepository) SetPasswordResetToken(ctx context.Context, id uint, token string, expiry *time.Time) error {
	return r.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_reset_token":  token,
			"password_reset_expiry": expiry,
		}).Error
}

// ClearPasswordResetToken 清除密码重置令牌
func (r *userRepository) ClearPasswordResetToken(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_reset_token":  nil,
			"password_reset_expiry": nil,
		}).Error
}

// Activate 激活用户
func (r *userRepository) Activate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Update("is_active", true).Error
}

// Deactivate 停用用户
func (r *userRepository) Deactivate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error
} 