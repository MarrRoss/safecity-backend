package zone

import (
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
)

type GetZonesByUserIDHandler struct {
	qryHandler *zone.GetZonesByUserIDHandler
	observer   *observability.Observability
}

func NewGetZonesByUserIDHandler(
	qryHandler *zone.GetZonesByUserIDHandler,
	observer *observability.Observability,
) *GetZonesByUserIDHandler {
	return &GetZonesByUserIDHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

func (h *GetZonesByUserIDHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.GetZonesByUserIDRequest
	//err := ctx.ParamsParser(&req)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
	//	resp := response.NewErrorResponse("failed to parse request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//
	//qry := zone.GetZonesByUserIDQuery{ID: req.UserID}
	//zones, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get user zones")
	//	return http.RespondError(ctx, err)
	//}
	//
	//zonesResp := response.NewGetZonesResponse(zones)
	//resp := response.NewSuccessResponse[[]*response.GetZoneResponse](zonesResp)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
