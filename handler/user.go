package handler

import (
	"errors"
	"net/http"

	models "github.com/MDmitryM/food_delivery_registration"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetUserByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int32)
	if !ok {
		return SendErrorJSON(ctx, http.StatusUnauthorized, errors.New("invalid userID"))
	}

	user, err := h.service.GetUserByID(ctx.Context(), userID)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusUnauthorized, err)
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (h *Handler) ChangeUserPassword(ctx *fiber.Ctx) error {
	var userPwdUpd models.UpdateUser

	userID, ok := ctx.Locals("user_id").(int32)
	if !ok {
		return SendErrorJSON(ctx, http.StatusUnauthorized, errors.New("invalid userID"))
	}

	if err := ctx.BodyParser(&userPwdUpd); err != nil {
		return SendErrorJSON(ctx, http.StatusBadRequest, err)
	}

	userPwdUpd.ID = userID

	user, err := h.service.UpdateUserPwd(ctx.Context(), userPwdUpd)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusInternalServerError, err)
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

type DeleteUaweResponce struct {
	Status string `json:"status"`
}

func (h *Handler) DeleteUserByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int32)
	if !ok {
		return SendErrorJSON(ctx, http.StatusUnauthorized, errors.New("invalid userID"))
	}

	rowsAffected, err := h.service.DeleteUserByID(ctx.Context(), userID)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusInternalServerError, err)
	}
	if rowsAffected == 0 {
		return SendErrorJSON(ctx, http.StatusNotFound, err)
	}
	return ctx.Status(http.StatusOK).JSON(DeleteUaweResponce{"ok"})
}
