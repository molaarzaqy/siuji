package middleware

import (
	"strings"

	"siuji-backend/pkg/jwt"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

func extractTempToken(c fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return ""
}

func TempAuth(jwtManager *jwt.Manager, log *logrus.Logger, expectedPurpose string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := extractTempToken(c)
		if tokenString == "" {
			log.Warnf("[TEMP_AUTH] missing temp token - ip: %s, path: %s", c.IP(), c.Path())
			return fiber.NewError(fiber.StatusUnauthorized, "missing temp token in Authorization header")
		}

		claims, err := jwtManager.ValidateTempToken(tokenString, expectedPurpose)
		if err != nil {
			log.Warnf("[TEMP_AUTH] invalid temp token - ip: %s, path: %s, error: %v", c.IP(), c.Path(), err)
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired temp token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("token_purpose", claims.Purpose)

		return c.Next()
	}
}