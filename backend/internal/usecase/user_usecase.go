package usecase

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.UserRepository
}

func NewUserUseCase(log *logrus.Logger, userRepository repository.UserRepository) *UserUseCase {
	return &UserUseCase{Log: log, UserRepository: userRepository}
}

func (c *UserUseCase) GetAll(filter, sort string, limit, offset int) ([]model.UserResponse, int64, error) {
	users, total, err := c.UserRepository.FindAllPagination(filter, sort, limit, offset)
	if err != nil {
		c.Log.Errorf("failed to fetch users: %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch users")
	}

	responses := make([]model.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *converter.UserToResponse(&u))
	}
	return responses, total, nil
}

func (c *UserUseCase) GetDetail(publicID string) (*model.UserResponse, error) {
	user, err := c.UserRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Delete(publicID string) error {
	if err := c.UserRepository.Delete(publicID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return nil
}