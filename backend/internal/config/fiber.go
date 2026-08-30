package config

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "siuji",
		ErrorHandler: NewErrorHandler(),
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