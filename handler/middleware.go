package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	authHeader = "Authorization"
	userCtx    = "user_id"
)

func (h *Handler) CheckToken(ctx *fiber.Ctx) error {
	header := ctx.Get(authHeader)
	if header == "" {
		return SendErrorJSON(ctx, http.StatusUnauthorized, errors.New("authorization is required"))
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return SendErrorJSON(ctx, http.StatusUnauthorized, errors.New("invalid auth header"))
	}

	token := headerParts[1]
	userID, err := h.service.ParseToken(token)
	if err != nil {
		return SendErrorJSON(ctx, http.StatusUnauthorized, err)
	}

	ctx.Locals(userCtx, userID)

	return ctx.Next()
}
