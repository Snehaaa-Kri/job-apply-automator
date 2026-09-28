package middleware

import (
	"strings"

	"github.com/Snehaaa-Kri/job-apply-automator/config"
	"github.com/Snehaaa-Kri/job-apply-automator/pkg/response"
	"github.com/Snehaaa-Kri/job-apply-automator/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired is a Gin middleware that validates the JWT token.
// Usage on a route: router.GET("/me", middleware.AuthRequired(cfg), handler)

func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	// gin.HandlerFunc = func(*gin.Context) — same as Express's (req, res, next) => void
	return func(c *gin.Context) {
		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort() // Stop processing — don't call the actual handler
			return
		}

		// 2. Header must be "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format, use: Bearer <token>")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// 3. Parse and validate the JWT
		// Node analogy: jwt.verify(token, secret)
		claims := &service.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what we expect (HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// 4. Attach user info to the context for downstream handlers
		// Node analogy: req.user = decoded
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// 5. Continue to the actual handler
		// Node analogy: next()
		c.Next()
	}
}
