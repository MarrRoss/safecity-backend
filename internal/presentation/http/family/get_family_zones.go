package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetFamilyZonesHandler struct {
	qryHandler *family.GetFamilyZonesHandler
	observer   *observability.Observability
}

func NewGetFamilyZonesHandler(
	qryHandler *family.GetFamilyZonesHandler,
	observer *observability.Observability,
) *GetFamilyZonesHandler {
	return &GetFamilyZonesHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get family zones
//	@Tags    zones
//	@Accept    json
//	@Produce  json
//	@Param     params      path   GetFamilyZonesRequest      true  "Search family zones"
//	@Success  200        {object}  response.BaseResponseData[[]response.GetZoneResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /families/{id}/zones [get]
//	@Security		APIKeyAuth
func (h *GetFamilyZonesHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetFamilyZonesRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	qry := family.GetFamilyZonesQuery{ID: req.ID}
	zones, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get family zones")
		return http.RespondError(ctx, err)
	}

	zonesResp := response.NewGetZonesResponse(zones)
	resp := response.NewSuccessResponse[[]*response.GetZoneResponse](zonesResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
