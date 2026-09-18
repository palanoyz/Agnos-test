package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ClaimsKey = "auth_claims"

func Middleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := ParseToken(secret, strings.TrimSpace(strings.TrimPrefix(header, prefix)))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}
