package cookie

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// SetAuthCookies sets the access and refresh token cookies.
// Durations are passed explicitly by the caller (e.g. from jwt.Manager),
// so this package doesn't depend on any global config.
func SetAuthCookies(c fiber.Ctx, accessToken, refreshToken string, accessDuration, refreshDuration time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false, // TODO: set true di production (bisa dari env APP_ENV)
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   int(accessDuration.Seconds()),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/api/v1/auth/refresh",
		MaxAge:   int(refreshDuration.Seconds()),
	})
}

func ClearAccessTokenCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
	})
}

func ClearAuthCookies(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Path:     "/api/v1/auth/refresh",
		MaxAge:   -1,
	})
}