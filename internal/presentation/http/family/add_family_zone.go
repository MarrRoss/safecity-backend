package family

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddFamilyZoneHandler struct {
	cmdHandler *family.AddFamilyZoneHandler
	observer   *observability.Observability
}

func NewAddFamilyZoneHandler(
	cmdHandler *family.AddFamilyZoneHandler,
	observer *observability.Observability,
) *AddFamilyZoneHandler {
	return &AddFamilyZoneHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add family zone
//	@Tags    zones
//	@Accept    json
//	@Produce  json
//	@Param   params   path   GetFamilyByIDRequest   true   "Search family"
//	@Param     request  body   AddZoneRequest  true  "Family and zone data"
//	@Success  200        {object}  response.BaseResponseData[response.AddFamilyZoneResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /families/{id}/zones [post]
//	@Security		APIKeyAuth
func (h *AddFamilyZoneHandler) Handle(ctx *fiber.Ctx) error {
	var pathReq request.GetFamilyByIDRequest
	if err := ctx.ParamsParser(&pathReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse id request")
		resp := response.NewErrorResponse("failed to parse id request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	var bodyReq request.AddZoneRequest
	if err := ctx.BodyParser(&bodyReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse body request")
		resp := response.NewErrorResponse("failed to parse body request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	boundaries := make([]app_model.ApplicationLocation, len(bodyReq.Boundaries))
	for i, b := range bodyReq.Boundaries {
		boundaries[i] = app_model.ApplicationLocation{
			Latitude:  b.Latitude,
			Longitude: b.Longitude,
		}
	}

	cmd := family.AddFamilyZoneCommand{
		FamilyID:   pathReq.ID,
		Name:       bodyReq.Name,
		Boundaries: boundaries,
		Safety:     bodyReq.Safety,
	}

	id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to create zone")
		return http.RespondError(ctx, err)
	}

	idResp := response.NewAddFamilyZoneResponse(id)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
