package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
)

type DeleteFamilyZoneHandler struct {
	cmdHandler *family.DeleteFamilyZoneHandler
	observer   *observability.Observability
}

func NewDeleteFamilyZoneHandler(
	cmdHandler *family.DeleteFamilyZoneHandler,
	observer *observability.Observability,
) *DeleteFamilyZoneHandler {
	return &DeleteFamilyZoneHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

func (h *DeleteFamilyZoneHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.DeleteFamilyZoneRequest
	//err := ctx.ParamsParser(&req)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
	//	resp := response.NewErrorResponse("failed to parse request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//
	//cmd := family.DeleteFamilyZoneCommand{
	//	FamilyID: req.FamilyID,
	//	ZoneID:   req.ZoneID,
	//}
	//err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to delete zone from family")
	//	return http.RespondError(ctx, err)
	//}
	//resp := response.NewSuccessResponse[any](nil)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
