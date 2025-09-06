package family

import (
	"awesomeProjectDDD/internal/application/family"
	"github.com/gofiber/fiber/v2"
)

type DeleteFamilyZonesHandler struct {
	cmdHandler *family.DeleteFamilyZonesHandler
}

func NewDeleteFamilyZonesHandler(
	cmdHandler *family.DeleteFamilyZonesHandler) *DeleteFamilyZonesHandler {
	return &DeleteFamilyZonesHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *DeleteFamilyZonesHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.DeleteFamilyZonesRequest
	//err := ctx.ParamsParser(&req)
	//if err != nil {
	//	return ctx.Status(fiber.StatusBadRequest).
	//		JSON(map[string]string{"status": "error", "detail": "error parsing request body"})
	//}
	//
	//cmd := family.DeleteFamilyZonesCommand{
	//	FamilyID: req.ID,
	//}
	//err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	//if err != nil {
	//	return ctx.Status(fiber.StatusInternalServerError).
	//		JSON(map[string]string{"status": "error", "detail": err.Error()})
	//}
	//return ctx.Status(fiber.StatusOK).
	//	JSON(map[string]string{"status": "ok"})
	return nil
}
