package middleware

import "github.com/gofiber/fiber/v3"

func RequireRole(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "role not found in token")
		}
		for _, allowedRole := range roles {
			if role == allowedRole {
				return c.Next()
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "you do not have permission to access this resource")
	}
}