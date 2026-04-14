package middleware

import (
	github.com/gin-gonic/gin
	github.com/golang-jwt/jwt/v5
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	time      time   time
)

var JWTSecret = []byte(your_jwt_secret)

type Claims struct {
	UserID string `json: user_id`
	Role   string `json: role`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
		authHeader := c.GetHeader(Authorization)
		if authHeader ==  {
			response.Error(c, 401, 未登录)
			c.Abort()
			return
		}

		tokenString := authHeader[len(Bearer ):]
		claims, err := ParseToken(tokenString)
		if err != nil {
			response.Error(c, 401, Token失效)
			c.Abort()
			return
		}

		c.Set(userID, claims.UserID)
		c.Set(role, claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(role)
		if !exists || role != admin {
			response.Error(c, 403, 无权限)
			c.Abort()
			return
		}
		c.Next()
	}
}
