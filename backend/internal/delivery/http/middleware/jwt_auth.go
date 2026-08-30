package middleware

import (
	"strings"

	"siuji-backend/pkg/cookie"
	"siuji-backend/pkg/jwt"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

func extractToken(c fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return c.Cookies("access_token")
}

func JWTAuth(jwtManager *jwt.Manager, log *logrus.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := extractToken(c)
		if tokenString == "" {
			log.Warnf("[AUTH] missing token - ip: %s, path: %s", c.IP(), c.Path())
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			log.Warnf("[AUTH] invalid token - ip: %s, error: %v", c.IP(), err)
			cookie.ClearAccessTokenCookie(c)
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("public_id", claims.PublicID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}