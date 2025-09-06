package role

import (
	"awesomeProjectDDD/internal/application/role"
	"github.com/gofiber/fiber/v2"
)

type AddRoleHandler struct {
	cmdHandler *role.AddRoleHandler
}

func NewAddRoleHandler(cmdHandler *role.AddRoleHandler) *AddRoleHandler {
	return &AddRoleHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *AddRoleHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.AddRoleRequest
	//err := ctx.BodyParser(&req)
	//if err != nil {
	//	return ctx.Status(fiber.StatusBadRequest).
	//		JSON(map[string]string{"status": "error", "detail": "error parsing request body"})
	//}
	//
	//cmd := role.AddRoleCommand{
	//	Name: req.Name,
	//}
	//
	//id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	//if err != nil {
	//	return ctx.Status(fiber.StatusInternalServerError).
	//		JSON(map[string]string{"status": "error", "detail": err.Error()})
	//}
	//return ctx.Status(fiber.StatusOK).
	//	JSON(map[string]string{"status": "ok", "detail": id})
	return nil
}
