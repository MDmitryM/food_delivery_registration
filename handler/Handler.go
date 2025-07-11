package handler

import (
	"net/http"

	"github.com/MDmitryM/food_delivery_registration/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Get("/sign-in", h.SignIn)
	app.Get("/sign-up", h.SignUp)
}

func (h *Handler) SignIn(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("SignIn")
}

func (h *Handler) SignUp(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("SignUp")
}
