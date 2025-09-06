package zone

import (
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
)

type AddZoneHandler struct {
	cmdHandler *zone.AddZoneHandler
	observer   *observability.Observability
}

func NewAddZoneHandler(
	cmdHandler *zone.AddZoneHandler,
	observer *observability.Observability,
) *AddZoneHandler {
	return &AddZoneHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

func (h *AddZoneHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.AddZoneRequest
	//err := ctx.BodyParser(&req)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
	//	resp := response.NewErrorResponse("failed to parse request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//
	//boundaries := make([]app_model.ApplicationLocation, len(req.Boundaries))
	//for i, b := range req.Boundaries {
	//	boundaries[i] = app_model.ApplicationLocation{
	//		Latitude:  b.Latitude,
	//		Longitude: b.Longitude,
	//	}
	//}
	//
	//cmd := zone.AddZoneCommand{
	//	Name:       req.Name,
	//	Boundaries: boundaries,
	//	Safety:     req.Safety,
	//}
	//
	//id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to add user zone")
	//	return http.RespondError(ctx, err)
	//}
	//
	//idResp := response.NewAddZoneResponse(id)
	//resp := response.NewSuccessResponse(idResp)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
