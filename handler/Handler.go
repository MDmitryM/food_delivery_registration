package handler

import (
	"github.com/MDmitryM/food_delivery_registration/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Post("/sign-in", h.SignIn)
	app.Post("/sign-up", h.SignUp)

	authorized := app.Group("/user", h.CheckToken)
	authorized.Get("/user-details", h.GetUserByID)
	authorized.Put("/change-password", h.ChangeUserPassword)
	authorized.Delete("/delete-user", h.DeleteUserByID)
}
