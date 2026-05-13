package middleware

import (
	"fileserver/internal/auth/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing auth header",
				},
			)
			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token",
				},
			)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
