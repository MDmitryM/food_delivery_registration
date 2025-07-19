package handler

import (
	"net/http"

	models "github.com/MDmitryM/food_delivery_registration"
	"github.com/gofiber/fiber/v2"
)

type SignInResponce struct {
	Token string `json:"access_token" validate:"required"`
}

func (h *Handler) SignIn(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return SendErrorJSON(ctx, http.StatusBadRequest, err)
	}

	if err := validate.Struct(&user); err != nil {
		return SendErrorJSON(ctx, http.StatusBadRequest, err)
	}

	token, err := h.service.IsUserValid(ctx.Context(), user)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusInternalServerError, err)
	}

	return ctx.Status(http.StatusOK).JSON(SignInResponce{token})
}

type SignUpResponce struct {
	UserID int32  `json:"user_id"`
	Token  string `json:"access_token"`
}

func (h *Handler) SignUp(ctx *fiber.Ctx) error {
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

	return ctx.Status(http.StatusOK).JSON(SignUpResponce{UserID: userID, Token: token})
}

type ValidationResponce struct {
	UserID int32 `json:"user_id"`
}

func (h *Handler) ValidateToken(ctx *fiber.Ctx) error {
	var input SignInResponce

	if err := ctx.BodyParser(&input); err != nil {
		return SendErrorJSON(ctx, http.StatusBadRequest, err)
	}

	userID, err := h.service.ParseToken(input.Token)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusUnauthorized, err)
	}

	_, err = h.service.GetUserByID(ctx.Context(), userID)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusUnauthorized, err)
	}

	return ctx.Status(http.StatusOK).JSON(ValidationResponce{userID})
}
