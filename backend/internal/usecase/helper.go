package usecase

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

// mapNotFoundError returns a 404 when the repository error message ends with
// "not found" (our convention across repositories), a 400 for invalid id
// formats, and otherwise logs the real error and returns a 500 — so unexpected
// failures are never silently disguised as "not found".
func mapNotFoundError(log *logrus.Logger, err error, context string) error {
	msg := err.Error()
	if strings.HasSuffix(msg, "not found") {
		return fiber.NewError(fiber.StatusNotFound, msg)
	}
	if strings.Contains(msg, "invalid uuid format") {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id format")
	}
	log.Errorf("%s: %+v", context, err)
	return fiber.NewError(fiber.StatusInternalServerError, context)
}