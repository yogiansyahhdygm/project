package middleware

import (
	"net/http"
	"strings"
	"toko/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ditemukan"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
			c.Set("claims", claims)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
			c.Abort()
			return
		}

	}
}

// hanya untuk admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses hanya untuk admin"})
			c.Abort()
			return
		}
		c.Next()
	}
}
