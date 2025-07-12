package handler

import (
	"net/http"

	models "github.com/MDmitryM/food_delivery_registration"
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
}

type SignInResponce struct {
	UserID int32  `json:"user_id"`
	Token  string `json:"access_token"`
}

func (h *Handler) SignIn(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return SendErrorJSON(ctx, http.StatusBadRequest, err)
	}

	if err := validate.Struct(&user); err != nil {
		return SendErrorJSON(ctx, http.StatusBadRequest, err)
	}

	userID, token, err := h.service.CreateUser(ctx.Context(), user)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusInternalServerError, err)
	}

	return ctx.Status(http.StatusOK).JSON(SignInResponce{UserID: userID, Token: token})
}

func (h *Handler) SignUp(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("SignUp")
}
