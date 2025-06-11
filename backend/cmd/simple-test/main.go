package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"temp-mailbox-service/internal/application"
	"temp-mailbox-service/internal/domain/user"
	"temp-mailbox-service/internal/infrastructure/auth"
)

// mockUserRepository 内存模拟用户仓储
type mockUserRepository struct {
	users  map[uint]*user.User
	emails map[string]*user.User
	usernames map[string]*user.User
	nextID uint
	mu     sync.RWMutex
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:     make(map[uint]*user.User),
		emails:    make(map[string]*user.User),
		usernames: make(map[string]*user.User),
		nextID:    1,
	}
}

func (r *mockUserRepository) Create(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	u.ID = r.nextID
	r.nextID++
	
	r.users[u.ID] = u
	r.emails[u.Email] = u
	r.usernames[u.Username] = u
	return nil
}

func (r *mockUserRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	u, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (r *mockUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	u, exists := r.emails[email]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (r *mockUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	u, exists := r.usernames[username]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (r *mockUserRepository) Update(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.users[u.ID]; !exists {
		return fmt.Errorf("用户不存在")
	}
	
	r.users[u.ID] = u
	return nil
}

func (r *mockUserRepository) Delete(ctx context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if u, exists := r.users[id]; exists {
		delete(r.users, id)
		delete(r.emails, u.Email)
		delete(r.usernames, u.Username)
	}
	return nil
}

func (r *mockUserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}

func (r *mockUserRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.users)), nil
}

func (r *mockUserRepository) Exists(ctx context.Context, id uint) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.users[id]
	return exists, nil
}

func (r *mockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.emails[email]
	return exists, nil
}

func (r *mockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.usernames[username]
	return exists, nil
}

func (r *mockUserRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if u, exists := r.users[id]; exists {
		u.Password = hashedPassword
		return nil
	}
	return fmt.Errorf("用户不存在")
}

func (r *mockUserRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	// 模拟实现，实际不更新
	return nil
}

func (r *mockUserRepository) SetPasswordResetToken(ctx context.Context, id uint, token string, expiry *time.Time) error {
	// 模拟实现
	return nil
}

func (r *mockUserRepository) ClearPasswordResetToken(ctx context.Context, id uint) error {
	// 模拟实现
	return nil
}

func (r *mockUserRepository) Activate(ctx context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if u, exists := r.users[id]; exists {
		u.IsActive = true
		return nil
	}
	return fmt.Errorf("用户不存在")
}

func (r *mockUserRepository) Deactivate(ctx context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if u, exists := r.users[id]; exists {
		u.IsActive = false
		return nil
	}
	return fmt.Errorf("用户不存在")
}

func main() {
	fmt.Println("🧪 开始测试用户系统核心组件（无数据库版本）...")

	// 1. 初始化模拟仓储和服务
	fmt.Println("⚙️  初始化模拟服务...")
	userRepo := newMockUserRepository()
	jwtService := auth.NewJWTService(
		"test-jwt-secret-key-for-testing-only-2024",
		60,    // 60分钟访问令牌
		10080, // 7天刷新令牌
		"temp-mailbox-test",
	)
	userService := application.NewUserService(userRepo, jwtService)
	fmt.Println("✅ 服务初始化完成")

	// 2. 测试用户注册
	fmt.Println("👤 测试用户注册...")
	ctx := context.Background()
	registerReq := &user.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	registeredUser, err := userService.Register(ctx, registerReq)
	if err != nil {
		log.Fatalf("❌ 用户注册失败: %v", err)
	}
	fmt.Printf("✅ 用户注册成功，ID: %d, 用户名: %s\n", registeredUser.ID, registeredUser.Username)

	// 3. 测试重复注册（应该失败）
	fmt.Println("🔄 测试重复注册...")
	_, err = userService.Register(ctx, registerReq)
	if err != nil {
		fmt.Printf("✅ 重复注册正确被拒绝: %v\n", err)
	} else {
		log.Fatalf("❌ 重复注册应该失败")
	}

	// 4. 测试用户登录
	fmt.Println("🔐 测试用户登录...")
	loginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := userService.Login(ctx, loginReq)
	if err != nil {
		log.Fatalf("❌ 用户登录失败: %v", err)
	}
	fmt.Printf("✅ 用户登录成功，访问令牌长度: %d\n", len(loginResp.Tokens.AccessToken))

	// 5. 测试错误密码登录（应该失败）
	fmt.Println("🚫 测试错误密码登录...")
	wrongLoginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	_, err = userService.Login(ctx, wrongLoginReq)
	if err != nil {
		fmt.Printf("✅ 错误密码登录正确被拒绝: %v\n", err)
	} else {
		log.Fatalf("❌ 错误密码登录应该失败")
	}

	// 6. 测试JWT令牌验证
	fmt.Println("🔑 测试JWT令牌验证...")
	claims, err := jwtService.ValidateAccessToken(loginResp.Tokens.AccessToken)
	if err != nil {
		log.Fatalf("❌ JWT令牌验证失败: %v", err)
	}
	fmt.Printf("✅ JWT令牌验证成功，用户ID: %d, 邮箱: %s\n", claims.UserID, claims.Email)

	// 7. 测试获取用户资料
	fmt.Println("📋 测试获取用户资料...")
	profile, err := userService.GetProfile(ctx, registeredUser.ID)
	if err != nil {
		log.Fatalf("❌ 获取用户资料失败: %v", err)
	}
	fmt.Printf("✅ 获取用户资料成功，全名: %s %s\n", profile.FirstName, profile.LastName)

	// 8. 测试更新用户资料
	fmt.Println("✏️  测试更新用户资料...")
	updateReq := &user.UpdateUserRequest{
		FirstName: "Updated",
		LastName:  "Name",
		TimeZone:  "Asia/Shanghai",
		Language:  "zh-CN",
	}

	updatedProfile, err := userService.UpdateProfile(ctx, registeredUser.ID, updateReq)
	if err != nil {
		log.Fatalf("❌ 更新用户资料失败: %v", err)
	}
	fmt.Printf("✅ 更新用户资料成功，新全名: %s %s\n", updatedProfile.FirstName, updatedProfile.LastName)

	// 9. 测试令牌刷新
	fmt.Println("🔄 测试令牌刷新...")
	newTokens, err := userService.RefreshToken(ctx, loginResp.Tokens.RefreshToken)
	if err != nil {
		log.Fatalf("❌ 令牌刷新失败: %v", err)
	}
	fmt.Printf("✅ 令牌刷新成功，新访问令牌长度: %d\n", len(newTokens.AccessToken))

	// 10. 测试修改密码
	fmt.Println("🔑 测试修改密码...")
	changePasswordReq := &user.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword456",
	}
	err = userService.ChangePassword(ctx, registeredUser.ID, changePasswordReq)
	if err != nil {
		log.Fatalf("❌ 修改密码失败: %v", err)
	}
	fmt.Println("✅ 修改密码成功")

	// 11. 测试使用新密码登录
	fmt.Println("🔐 测试使用新密码登录...")
	newLoginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "newpassword456",
	}
	_, err = userService.Login(ctx, newLoginReq)
	if err != nil {
		log.Fatalf("❌ 新密码登录失败: %v", err)
	}
	fmt.Println("✅ 新密码登录成功")

	// 12. 测试用户列表
	fmt.Println("📋 测试用户列表...")
	users, total, err := userService.ListUsers(ctx, 0, 10)
	if err != nil {
		log.Fatalf("❌ 获取用户列表失败: %v", err)
	}
	fmt.Printf("✅ 获取用户列表成功，总数: %d, 当前页: %d个用户\n", total, len(users))

	fmt.Println("\n🎉 所有测试通过！用户系统核心组件运行正常！")
	fmt.Println("📊 测试总结:")
	fmt.Println("   ✅ 用户注册和唯一性验证")
	fmt.Println("   ✅ 用户登录和密码验证")
	fmt.Println("   ✅ JWT令牌生成和验证")
	fmt.Println("   ✅ 用户资料管理")
	fmt.Println("   ✅ 令牌刷新")
	fmt.Println("   ✅ 密码修改")
	fmt.Println("   ✅ 用户列表查询")
	fmt.Println("\n💡 注意：此测试使用内存模拟仓储，实际部署时需要连接真实数据库")
} 