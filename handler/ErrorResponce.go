package handler

import (
	"github.com/gofiber/fiber/v2"
)

type MyErrorResponce struct {
	Error string `json:"error"`
}

func SendErrorJSON(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(MyErrorResponce{err.Error()})
}
