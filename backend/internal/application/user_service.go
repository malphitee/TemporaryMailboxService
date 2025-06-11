package application

import (
	"context"
	"errors"
	"fmt"

	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/auth"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	// 认证相关
	Register(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error)
	Login(ctx context.Context, req *user.LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error)
	
	// 用户管理
	GetProfile(ctx context.Context, userID uint) (*user.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *user.UpdateUserRequest) (*user.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *user.ChangePasswordRequest) error
	
	// 用户查询
	GetUserByID(ctx context.Context, userID uint) (*user.UserResponse, error)
	ListUsers(ctx context.Context, offset, limit int) ([]*user.UserResponse, int64, error)
	
	// 账户管理
	DeactivateUser(ctx context.Context, userID uint) error
	ActivateUser(ctx context.Context, userID uint) error
}

// LoginResponse 登录响应
type LoginResponse struct {
	User   *user.UserResponse `json:"user"`
	Tokens *auth.TokenPair    `json:"tokens"`
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

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}
	
	// 检查用户名是否已存在
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, errors.New("用户名已被占用")
	}
	
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}
	
	// 创建用户实体
	newUser := &user.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
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

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *user.LoginRequest) (*LoginResponse, error) {
	// 根据邮箱查找用户
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if existingUser == nil {
		return nil, errors.New("邮箱或密码错误")
	}
	
	// 检查用户是否已激活
	if !existingUser.IsActive {
		return nil, errors.New("账户已被停用")
	}
	
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}
	
	// 生成JWT令牌
	tokens, err := s.jwtService.GenerateTokens(existingUser.ID, existingUser.Username, existingUser.Email)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}
	
	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, existingUser.ID); err != nil {
		// 记录错误但不影响登录流程
		// TODO: 添加日志记录
	}
	
	return &LoginResponse{
		User:   existingUser.ToResponse(),
		Tokens: tokens,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	return s.jwtService.RefreshTokens(refreshToken)
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID uint) (*user.UserResponse, error) {
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if existingUser == nil {
		return nil, errors.New("用户不存在")
	}
	
	return existingUser.ToResponse(), nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	// 获取现有用户
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if existingUser == nil {
		return nil, errors.New("用户不存在")
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
	// 获取现有用户
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if existingUser == nil {
		return errors.New("用户不存在")
	}
	
	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("当前密码错误")
	}
	
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	
	// 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}
	
	return nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(ctx context.Context, userID uint) (*user.UserResponse, error) {
	return s.GetProfile(ctx, userID)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]*user.UserResponse, int64, error) {
	// 获取用户列表
	users, err := s.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}
	
	// 获取总数
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户总数失败: %w", err)
	}
	
	// 转换为响应格式
	responses := make([]*user.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}
	
	return responses, total, nil
}

// DeactivateUser 停用用户
func (s *userService) DeactivateUser(ctx context.Context, userID uint) error {
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("检查用户失败: %w", err)
	}
	if !exists {
		return errors.New("用户不存在")
	}
	
	return s.userRepo.Deactivate(ctx, userID)
}

// ActivateUser 激活用户
func (s *userService) ActivateUser(ctx context.Context, userID uint) error {
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("检查用户失败: %w", err)
	}
	if !exists {
		return errors.New("用户不存在")
	}
	
	return s.userRepo.Activate(ctx, userID)
} 