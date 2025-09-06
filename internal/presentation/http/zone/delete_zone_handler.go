package zone

import (
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type DeleteZoneHandler struct {
	qryHandler *zone.DeleteZoneHandler
	observer   *observability.Observability
}

func NewDeleteZoneHandler(
	qryHandler *zone.DeleteZoneHandler,
	observer *observability.Observability) *DeleteZoneHandler {
	return &DeleteZoneHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Delete zone
//	@Tags    zones
//	@Accept    json
//	@Produce  json
//	@Param     params      path   DeleteZoneRequest      true  "Delete family zone"
//	@Success  200        {object}  response.BaseResponseData[string]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /zones/{id} [delete]
//	@Security		APIKeyAuth
func (h *DeleteZoneHandler) Handle(ctx *fiber.Ctx) error {
	var req request.DeleteZoneRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "enable to parse request"})
	}
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	qry := zone.DeleteZoneQuery{
		ID:     req.ID,
		UserID: *identity.Sub}
	err = h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to delete zone")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
