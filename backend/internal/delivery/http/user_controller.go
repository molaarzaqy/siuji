package http

import (
	"strconv"

	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase) *UserController {
	return &UserController{UseCase: useCase}
}

func (ctrl *UserController) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	filter := c.Query("filter", "")
	sort := c.Query("sort", "")
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, total, err := ctrl.UseCase.GetAll(filter, sort, limit, offset)
	if err != nil {
		return err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := response.PaginationMeta{Page: page, Limit: limit, TotalDatas: int(total), TotalPages: totalPages, Filter: filter, Sort: sort}
	return response.SuccessPagination(c, "List user retrieved successfully.", users, meta)
}

func (ctrl *UserController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("user_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "User detail retrieved successfully.", result)
}

func (ctrl *UserController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("user_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "User deleted successfully.")
}