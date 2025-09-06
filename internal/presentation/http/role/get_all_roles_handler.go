package role

import (
	"awesomeProjectDDD/internal/application/role"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAllRolesHandler struct {
	qryHandler *role.GetAllRolesHandler
}

func NewGetAllRolesHandler(qryHandler *role.GetAllRolesHandler) *GetAllRolesHandler {
	return &GetAllRolesHandler{
		qryHandler: qryHandler,
	}
}

func (h *GetAllRolesHandler) Handle(ctx *fiber.Ctx) error {
	roles, err := h.qryHandler.Handle(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	resp := response.NewGetRolesResponse(roles)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
