package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"demo/internal/services"
)

const UserIDKey = "userID"

func Auth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "en-tête Authorization manquant"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "en-tête Authorization malformé, attendu: Bearer <token>"})
			return
		}

		userID, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide ou expiré"})
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
