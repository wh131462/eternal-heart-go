package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig 配置项结构体
type JWTConfig struct {
	SigningKey      string        // 必填，签名密钥
	AuthScheme      string        // 认证头前缀，默认"Bearer"
	ContextKey      string        // 用户信息在context中的键名
	TokenLookup     string        // Token获取方式，默认"header:Authorization"
	Expires         time.Duration // Token有效期
	RefreshInterval time.Duration // Token刷新间隔
}

// 默认配置
var DefaultConfig = JWTConfig{
	AuthScheme:      "Bearer",
	ContextKey:      "jwtClaims",
	TokenLookup:     "header:Authorization",
	Expires:         24 * time.Hour,
	RefreshInterval: 1 * time.Hour,
}

// 自定义错误类型
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("malformed token")
	ErrTokenInvalid     = errors.New("invalid token")
)

// 用户声明结构体
type CustomClaims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"uname"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 中间件入口
func JWTAuth(config JWTConfig) gin.HandlerFunc {
	// 合并配置
	conf := mergeConfig(config)

	return func(c *gin.Context) {
		// 白名单路径检查
		if skipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 解析Token
		tokenString, err := extractToken(c, conf)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// 验证并解析Token
		claims, err := parseToken(tokenString, conf)
		if err != nil {
			handleTokenError(c, err)
			return
		}

		// 注入用户信息到上下文
		c.Set(conf.ContextKey, claims)

		// Token自动刷新机制
		if shouldRefresh(claims) {
			newToken, _ := refreshToken(claims, conf)
			c.Header("X-New-Token", newToken)
		}

		c.Next()
	}
}

func shouldRefresh(claims *CustomClaims) bool {
	return true
}

// --- 辅助函数实现 ---

// 合并配置
func mergeConfig(conf JWTConfig) JWTConfig {
	if conf.AuthScheme == "" {
		conf.AuthScheme = DefaultConfig.AuthScheme
	}
	// 其他字段合并逻辑...
	return conf
}

// 提取Token
func extractToken(c *gin.Context, conf JWTConfig) (string, error) {
	// 支持多种获取方式：header/form/query
	parts := strings.Split(conf.TokenLookup, ":")
	switch parts[0] {
	case "header":
		return extractHeaderToken(c, parts[1], conf.AuthScheme)
	case "query":
		// 实现query参数解析
	case "form":
		// 实现form参数解析
	}
	return "", ErrTokenInvalid
}

// 解析头部的Token
func extractHeaderToken(c *gin.Context, headerName, authScheme string) (string, error) {
	authHeader := c.GetHeader(headerName)
	if authHeader == "" {
		return "", ErrTokenInvalid
	}

	// 分割Bearer和实际Token
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == authScheme) {
		return "", ErrTokenMalformed
	}
	return parts[1], nil
}

// 解析验证Token
func parseToken(tokenString string, conf JWTConfig) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(conf.SigningKey), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	// 错误类型细化
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, ErrTokenMalformed
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrTokenExpired
	} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, ErrTokenNotValidYet
	}
	return nil, ErrTokenInvalid
}

// Token刷新逻辑
func refreshToken(claims *CustomClaims, conf JWTConfig) (string, error) {
	// 延长有效期
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(conf.Expires))
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(conf.SigningKey))
}

// 错误处理
func handleTokenError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrTokenExpired):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40101, "error": "token expired"})
	case errors.Is(err, ErrTokenNotValidYet):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40102, "error": "token not active"})
	case errors.Is(err, ErrTokenMalformed):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40103, "error": "invalid token format"})
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40100, "error": "authentication failed"})
	}
}

// 白名单路径检查
func skipAuth(path string) bool {
	excludePaths := map[string]bool{
		"/api/login":    true,
		"/api/register": true,
		"/health":       true,
	}
	return excludePaths[path]
}
