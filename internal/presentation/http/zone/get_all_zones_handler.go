package zone

import (
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAllZonesHandler struct {
	cmdHandler *zone.GetAllZonesHandler
}

func NewGetAllZonesHandler(cmdHandler *zone.GetAllZonesHandler) *GetAllZonesHandler {
	return &GetAllZonesHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *GetAllZonesHandler) Handle(ctx *fiber.Ctx) error {
	zones, err := h.cmdHandler.Handle(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}

	resp := response.NewGetZonesResponse(zones)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
