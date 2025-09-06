package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/presentation/http/request"
	"github.com/gofiber/fiber/v2"
)

type RestoreUserHandler struct {
	qryHandler *user.RestoreUserHandler
}

func NewRestoreUserHandler(qryHandler *user.RestoreUserHandler) *RestoreUserHandler {
	return &RestoreUserHandler{
		qryHandler: qryHandler,
	}
}

func (h *RestoreUserHandler) Handle(ctx *fiber.Ctx) error {
	var req request.RestoreUserRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "enable to parse request"})
	}

	qry := user.RestoreUserQuery{ID: req.ID}
	err = h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	return ctx.
		Status(fiber.StatusOK).
		JSON(map[string]string{"status": "ok", "detail": "user successfully restored"})
}
