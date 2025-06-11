package application

import (
	"context"
	"fmt"

	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/auth"
)

// UserService 用户服务接口
type UserService interface {
	// 认证相关
	RegisterUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error)
	LoginUser(ctx context.Context, req *user.LoginRequest) (*LoginResponse, error)
	
	// 用户管理
	GetUserProfile(ctx context.Context, userID uint) (*user.UserResponse, error)
	UpdateUserProfile(ctx context.Context, userID uint, req *user.UpdateUserRequest) (*user.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *user.ChangePasswordRequest) error
}

// LoginResponse 登录响应
type LoginResponse struct {
	User   *user.UserResponse `json:"user"`
	Token  *auth.TokenPair    `json:"token"`
}

// userService 用户服务实现
type userService struct {
	userRepo   user.Repository
	jwtService auth.JWTService
}

// NewUserService 创建新的用户服务实例
func NewUserService(userRepo user.Repository, jwtService auth.JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// RegisterUser 用户注册
func (s *userService) RegisterUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("邮箱已存在")
	}
	
	// 检查用户名是否已存在
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("用户名已存在")
	}
	
	// 验证密码强度
	if err := auth.IsValidPassword(req.Password); err != nil {
		return nil, fmt.Errorf("密码不符合要求: %w", err)
	}
	
	// 哈希密码
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}
	
	// 创建用户实体
	newUser := &user.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}
	
	// 保存用户
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}
	
	return newUser.ToResponse(), nil
}

// LoginUser 用户登录
func (s *userService) LoginUser(ctx context.Context, req *user.LoginRequest) (*LoginResponse, error) {
	// 根据邮箱获取用户
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}
	
	// 检查用户是否激活
	if !existingUser.IsActive {
		return nil, fmt.Errorf("用户账户已被停用")
	}
	
	// 验证密码
	if !auth.VerifyPassword(existingUser.Password, req.Password) {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}
	
	// 生成JWT令牌
	tokenPair, err := s.jwtService.GenerateTokens(existingUser.ID, existingUser.Username, existingUser.Email)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}
	
	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, existingUser.ID); err != nil {
		// 记录日志但不影响登录流程
		fmt.Printf("更新最后登录时间失败: %v\n", err)
	}
	
	return &LoginResponse{
		User:   existingUser.ToResponse(),
		Token:  tokenPair,
	}, nil
}

// GetUserProfile 获取用户资料
func (s *userService) GetUserProfile(ctx context.Context, userID uint) (*user.UserResponse, error) {
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	
	return existingUser.ToResponse(), nil
}

// UpdateUserProfile 更新用户资料
func (s *userService) UpdateUserProfile(ctx context.Context, userID uint, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	
	// 更新用户信息
	existingUser.FirstName = req.FirstName
	existingUser.LastName = req.LastName
	existingUser.Avatar = req.Avatar
	existingUser.TimeZone = req.TimeZone
	existingUser.Language = req.Language
	
	// 保存更新
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}
	
	return existingUser.ToResponse(), nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, req *user.ChangePasswordRequest) error {
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	if existingUser == nil {
		return fmt.Errorf("用户不存在")
	}
	
	// 验证当前密码
	if !auth.VerifyPassword(existingUser.Password, req.CurrentPassword) {
		return fmt.Errorf("当前密码错误")
	}
	
	// 验证新密码强度
	if err := auth.IsValidPassword(req.NewPassword); err != nil {
		return fmt.Errorf("新密码不符合要求: %w", err)
	}
	
	// 哈希新密码
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}
	
	// 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}
	
	return nil
} 