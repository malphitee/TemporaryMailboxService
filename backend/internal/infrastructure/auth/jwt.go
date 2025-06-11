package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService JWT服务接口
type JWTService interface {
	GenerateTokens(userID uint, username, email string) (*TokenPair, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshToken(tokenString string) (*JWTClaims, error)
	RefreshTokens(refreshToken string) (*TokenPair, error)
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// jwtService JWT服务实现
type jwtService struct {
	secretKey           []byte
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
	issuer              string
}

// NewJWTService 创建新的JWT服务实例
func NewJWTService(secretKey string, accessTokenTTL, refreshTokenTTL int, issuer string) JWTService {
	return &jwtService{
		secretKey:           []byte(secretKey),
		accessTokenTTL:      time.Duration(accessTokenTTL) * time.Minute,
		refreshTokenTTL:     time.Duration(refreshTokenTTL) * time.Minute,
		issuer:              issuer,
	}
}

// GenerateTokens 生成访问令牌和刷新令牌
func (s *jwtService) GenerateTokens(userID uint, username, email string) (*TokenPair, error) {
	now := time.Now()
	
	// 生成访问令牌
	accessToken, err := s.generateToken(userID, username, email, now.Add(s.accessTokenTTL), "access")
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}
	
	// 生成刷新令牌
	refreshToken, err := s.generateToken(userID, username, email, now.Add(s.refreshTokenTTL), "refresh")
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}, nil
}

// ValidateAccessToken 验证访问令牌
func (s *jwtService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString, "access")
}

// ValidateRefreshToken 验证刷新令牌
func (s *jwtService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString, "refresh")
}

// RefreshTokens 使用刷新令牌获取新的令牌对
func (s *jwtService) RefreshTokens(refreshToken string) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("刷新令牌无效: %w", err)
	}
	
	// 生成新的令牌对
	return s.GenerateTokens(claims.UserID, claims.Username, claims.Email)
}

// generateToken 生成JWT令牌
func (s *jwtService) generateToken(userID uint, username, email string, expiresAt time.Time, tokenType string) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("user_%d", userID),
			ID:        fmt.Sprintf("%s_%d_%d", tokenType, userID, time.Now().Unix()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// validateToken 验证JWT令牌
func (s *jwtService) validateToken(tokenString, expectedType string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("令牌解析失败: %w", err)
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("令牌无效")
	}
	
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("无效的令牌声明")
	}
	
	// 验证令牌类型（通过ID前缀）
	if expectedType != "" && claims.ID != "" {
		expectedPrefix := fmt.Sprintf("%s_%d_", expectedType, claims.UserID)
		if len(claims.ID) < len(expectedPrefix) || claims.ID[:len(expectedPrefix)] != expectedPrefix {
			return nil, fmt.Errorf("令牌类型不匹配")
		}
	}
	
	return claims, nil
}

// ExtractTokenFromHeader 从Authorization头部提取令牌
func ExtractTokenFromHeader(authHeader string) string {
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	return ""
} 