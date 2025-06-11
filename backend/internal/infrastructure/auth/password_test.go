package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "有效密码",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "最小长度密码",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "长密码",
			password: "this_is_a_very_long_password_with_many_characters_123456789",
			wantErr:  false,
		},
		{
			name:     "密码太短",
			password: "12345",
			wantErr:  true,
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash, "哈希后的密码不应该等于原密码")
				assert.True(t, strings.HasPrefix(hash, "$2a$"), "应该使用bcrypt算法")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// 准备测试数据
	password := "testpassword123"
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		want           bool
	}{
		{
			name:           "正确密码",
			hashedPassword: hash,
			password:       password,
			want:           true,
		},
		{
			name:           "错误密码",
			hashedPassword: hash,
			password:       "wrongpassword",
			want:           false,
		},
		{
			name:           "空密码",
			hashedPassword: hash,
			password:       "",
			want:           false,
		},
		{
			name:           "无效哈希",
			hashedPassword: "invalid_hash",
			password:       password,
			want:           false,
		},
		{
			name:           "空哈希",
			hashedPassword: "",
			password:       password,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyPassword(tt.hashedPassword, tt.password)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "有效密码",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "最小长度密码",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "包含特殊字符的密码",
			password: "pass@word123!",
			wantErr:  false,
		},
		{
			name:     "密码太短",
			password: "12345",
			wantErr:  true,
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidPassword(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordConstants(t *testing.T) {
	// 验证常量值
	assert.Equal(t, 6, MinPasswordLength, "最小密码长度应该是6")
	assert.Equal(t, 12, BcryptCost, "Bcrypt成本应该是12")
}

func TestHashPasswordConsistency(t *testing.T) {
	// 测试相同密码多次哈希产生不同结果（盐值不同）
	password := "testpassword"
	
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "相同密码的哈希结果应该不同（因为盐值不同）")
	
	// 但都应该能验证原密码
	assert.True(t, VerifyPassword(hash1, password))
	assert.True(t, VerifyPassword(hash2, password))
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"
	
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword(hash, password)
	}
} 