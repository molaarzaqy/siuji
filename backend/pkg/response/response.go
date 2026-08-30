package response

import "github.com/gofiber/fiber/v3"

type Response struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type ResponseNoData struct {
	Status       string `json:"status"`
	ResponseCode int    `json:"response_code"`
	Message      string `json:"message,omitempty"`
}

type ResponsePaginated struct {
	Status       string         `json:"status"`
	ResponseCode int            `json:"response_code"`
	Message      string         `json:"message,omitempty"`
	Data         interface{}    `json:"data,omitempty"`
	Meta         PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalDatas int    `json:"total_datas"`
	TotalPages int    `json:"total_pages"`
	Filter     string `json:"filter"`
	Sort       string `json:"sort"`
}

func Success(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status: "success", ResponseCode: fiber.StatusOK, Message: message, Data: data,
	})
}

func SuccessNoData(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(ResponseNoData{
		Status: "success", ResponseCode: fiber.StatusOK, Message: message,
	})
}

func SuccessPagination(c fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(ResponsePaginated{
		Status: "success", ResponseCode: fiber.StatusOK, Message: message, Data: data, Meta: meta,
	})
}

func Created(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status: "success", ResponseCode: fiber.StatusCreated, Message: message, Data: data,
	})
}