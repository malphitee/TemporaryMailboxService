package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength 最小密码长度
	MinPasswordLength = 6
	// BcryptCost bcrypt算法成本
	BcryptCost = 12
)

// HashPassword 对密码进行哈希
func HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", fmt.Errorf("密码长度不能少于%d个字符", MinPasswordLength)
	}
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}
	
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// IsValidPassword 检查密码强度
func IsValidPassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("密码长度不能少于%d个字符", MinPasswordLength)
	}
	
	// 可以在这里添加更多密码强度检查
	// 例如：必须包含数字、大小写字母、特殊字符等
	
	return nil
} 