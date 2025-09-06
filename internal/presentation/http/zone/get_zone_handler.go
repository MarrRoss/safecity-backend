package zone

import (
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetZoneHandler struct {
	qryHandler *zone.GetZoneHandler
}

func NewGetZoneHandler(qryHandler *zone.GetZoneHandler) *GetZoneHandler {
	return &GetZoneHandler{
		qryHandler: qryHandler,
	}
}

func (gh *GetZoneHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetZoneByIDRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "error parsing request"})
	}
	qry := zone.GetZoneQuery{ID: req.ID.String()}
	z, err := gh.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	resp := response.NewGetZoneResponse(z)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
