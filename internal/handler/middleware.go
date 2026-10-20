package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"manajemen-barang/internal/models"
)

const CookieName = "token"

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			respondModelError(c, models.ErrUnauthorized)
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			respondModelError(c, models.ErrUnauthorized)
			c.Abort()
			return
		}

		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if sub, ok := claims["sub"].(float64); ok {
			c.Set("user_id", int64(sub))
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if v, err := c.Cookie(CookieName); err == nil && v != "" {
		return v
	}
	header := c.GetHeader("Authorization")
	if len(header) > 7 && strings.EqualFold(header[:7], "Bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return ""
}
