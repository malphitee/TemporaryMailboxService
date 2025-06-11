package auth

import (
	"strings"
	"testing"
	"time"
)

func TestJWTService_TokenGeneration(t *testing.T) {
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		60,    // 1小时
		10080, // 7天
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 测试令牌对生成
	tokens, err := jwtService.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成令牌对失败: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Error("访问令牌不应该为空")
	}
	if tokens.RefreshToken == "" {
		t.Error("刷新令牌不应该为空")
	}
	if tokens.TokenType != "Bearer" {
		t.Errorf("期望令牌类型为 'Bearer'，得到 '%s'", tokens.TokenType)
	}
	if tokens.ExpiresIn <= 0 {
		t.Error("过期时间应该大于0")
	}
}

func TestJWTService_TokenValidation(t *testing.T) {
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		60,    // 1小时
		10080, // 7天
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 生成令牌对
	tokens, err := jwtService.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成令牌对失败: %v", err)
	}

	// 验证访问令牌
	claims, err := jwtService.ValidateAccessToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("验证访问令牌失败: %v", err)
	}

	// 验证claims内容
	if claims.UserID != userID {
		t.Errorf("期望用户ID为 %d，得到 %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("期望邮箱为 '%s'，得到 '%s'", email, claims.Email)
	}
	if claims.Username != username {
		t.Errorf("期望用户名为 '%s'，得到 '%s'", username, claims.Username)
	}
	if claims.Issuer != "test-issuer" {
		t.Errorf("期望发行者为 'test-issuer'，得到 '%s'", claims.Issuer)
	}
	
	// 验证令牌类型（通过ID前缀）
	expectedPrefix := "access_123_"
	if !strings.HasPrefix(claims.ID, expectedPrefix) {
		t.Errorf("期望访问令牌ID以 '%s' 开头，得到 '%s'", expectedPrefix, claims.ID)
	}
}

func TestJWTService_RefreshTokenValidation(t *testing.T) {
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		60,    // 1小时
		10080, // 7天
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 生成令牌对
	tokens, err := jwtService.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成令牌对失败: %v", err)
	}

	// 验证刷新令牌
	claims, err := jwtService.ValidateRefreshToken(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("验证刷新令牌失败: %v", err)
	}

	// 验证claims内容
	if claims.UserID != userID {
		t.Errorf("期望用户ID为 %d，得到 %d", userID, claims.UserID)
	}
	
	// 验证令牌类型（通过ID前缀）
	expectedPrefix := "refresh_123_"
	if !strings.HasPrefix(claims.ID, expectedPrefix) {
		t.Errorf("期望刷新令牌ID以 '%s' 开头，得到 '%s'", expectedPrefix, claims.ID)
	}
}

func TestJWTService_InvalidToken(t *testing.T) {
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		60,    // 1小时
		10080, // 7天
		"test-issuer",
	)

	// 测试无效令牌
	invalidToken := "invalid.token.here"
	_, err := jwtService.ValidateAccessToken(invalidToken)
	if err == nil {
		t.Error("无效令牌应该验证失败")
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// 创建一个过期时间很短的JWT服务
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		-1,    // 负数表示已过期
		10080, // 7天
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 生成过期令牌
	tokens, err := jwtService.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成令牌对失败: %v", err)
	}

	// 验证过期令牌应该失败
	_, err = jwtService.ValidateAccessToken(tokens.AccessToken)
	if err == nil {
		t.Error("过期令牌应该验证失败")
	}
}

func TestJWTService_WrongTokenType(t *testing.T) {
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		60,    // 1小时
		10080, // 7天
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 生成令牌对
	tokens, err := jwtService.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成令牌对失败: %v", err)
	}

	// 用刷新令牌验证方法验证访问令牌（应该失败）
	_, err = jwtService.ValidateRefreshToken(tokens.AccessToken)
	if err == nil {
		t.Error("用错误类型验证令牌应该失败")
	}
}

func TestJWTService_RefreshTokens(t *testing.T) {
	jwtService := NewJWTService(
		"test-secret-key-for-testing-only",
		60,    // 1小时
		10080, // 7天
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 生成原始令牌对
	originalTokens, err := jwtService.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成原始令牌对失败: %v", err)
	}

	// 等待一秒以确保时间戳不同
	time.Sleep(1 * time.Second)

	// 使用刷新令牌获取新令牌对
	newTokens, err := jwtService.RefreshTokens(originalTokens.RefreshToken)
	if err != nil {
		t.Fatalf("刷新令牌失败: %v", err)
	}

	// 验证新令牌对
	if newTokens.AccessToken == "" {
		t.Error("新访问令牌不应该为空")
	}
	if newTokens.RefreshToken == "" {
		t.Error("新刷新令牌不应该为空")
	}
	
	// 验证新访问令牌可以被验证
	newClaims, err := jwtService.ValidateAccessToken(newTokens.AccessToken)
	if err != nil {
		t.Errorf("新访问令牌验证失败: %v", err)
	}
	
	// 验证用户信息保持一致
	if newClaims.UserID != userID {
		t.Errorf("期望用户ID为 %d，得到 %d", userID, newClaims.UserID)
	}
	if newClaims.Email != email {
		t.Errorf("期望邮箱为 '%s'，得到 '%s'", email, newClaims.Email)
	}
	if newClaims.Username != username {
		t.Errorf("期望用户名为 '%s'，得到 '%s'", username, newClaims.Username)
	}
	
	// 验证令牌确实不同（或者至少功能正常）
	originalClaims, err := jwtService.ValidateAccessToken(originalTokens.AccessToken)
	if err != nil {
		t.Errorf("原始访问令牌验证失败: %v", err)
	}
	
	// 比较令牌ID（时间戳应该不同）
	if originalClaims.ID == newClaims.ID && newTokens.AccessToken == originalTokens.AccessToken {
		t.Error("刷新后的令牌与原始令牌完全相同，这可能不是期望的行为")
	}
}

func TestJWTService_DifferentSecrets(t *testing.T) {
	jwtService1 := NewJWTService(
		"secret-key-1",
		60,
		10080,
		"test-issuer",
	)

	jwtService2 := NewJWTService(
		"secret-key-2", // 不同的密钥
		60,
		10080,
		"test-issuer",
	)

	userID := uint(123)
	email := "test@example.com"
	username := "testuser"

	// 用第一个服务生成令牌
	tokens, err := jwtService1.GenerateTokens(userID, username, email)
	if err != nil {
		t.Fatalf("生成令牌对失败: %v", err)
	}

	// 用第二个服务验证令牌（不同密钥，应该失败）
	_, err = jwtService2.ValidateAccessToken(tokens.AccessToken)
	if err == nil {
		t.Error("用不同密钥验证令牌应该失败")
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	// 测试正确的Bearer格式
	authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	token := ExtractTokenFromHeader(authHeader)
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	if token != expected {
		t.Errorf("期望提取的令牌为 '%s'，得到 '%s'", expected, token)
	}

	// 测试错误格式
	wrongHeader := "Basic dXNlcjpwYXNzd29yZA=="
	token = ExtractTokenFromHeader(wrongHeader)
	if token != "" {
		t.Errorf("错误格式应该返回空字符串，得到 '%s'", token)
	}

	// 测试空头部
	token = ExtractTokenFromHeader("")
	if token != "" {
		t.Errorf("空头部应该返回空字符串，得到 '%s'", token)
	}
} 