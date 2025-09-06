package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetUserFamiliesHandler struct {
	cmdHandler *family.GetUserFamiliesHandler
}

func NewGetUserFamiliesHandler(cmdHandler *family.GetUserFamiliesHandler) *GetUserFamiliesHandler {
	return &GetUserFamiliesHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *GetUserFamiliesHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetUserFamiliesRequest
	req.UserID = ctx.Params("user_id")
	if req.UserID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "missing user_id"})
	}

	qry := family.GetUserFamiliesQuery{ID: req.UserID}
	families, err := h.cmdHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	resp := response.NewGetFamiliesResponse(families)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
