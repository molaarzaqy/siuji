package config

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "siuji",
		ErrorHandler: NewErrorHandler(),
	})

	allowedOrigins := config.GetString("CORS_ORIGIN")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:5173"
	}
	app.Use(func(c fiber.Ctx) error {
		origin := c.Get("Origin")
		allowed := false
		for _, configuredOrigin := range strings.Split(allowedOrigins, ",") {
			if strings.TrimSpace(configuredOrigin) == origin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "internal server error"

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			message = e.Message
		}

		return c.Status(code).JSON(fiber.Map{
			"status":        "error",
			"response_code": code,
			"message":       message,
		})
	}
}
