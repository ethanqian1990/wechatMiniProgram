package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/response"
	"strings"
	"time"
)

var JWTSecret = []byte("your_jwt_secret")
var JWTExpire = 168 * time.Hour

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func SetJWTSecret(secret string) {
	if strings.TrimSpace(secret) == "" {
		return
	}
	JWTSecret = []byte(secret)
}

func SetJWTExpire(d time.Duration) {
	if d <= 0 {
		return
	}
	JWTExpire = d
}

func GenerateToken(userID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "未登录")
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) || len(authHeader) <= len(prefix) {
			response.Error(c, 401, "Token格式错误")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(authHeader[len(prefix):])
		claims, err := ParseToken(tokenString)
		if err != nil {
			response.Error(c, 401, "Token失效")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			response.Error(c, 403, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
