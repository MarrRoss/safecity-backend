package zone

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type ChangeZoneHandler struct {
	cmdHandler *zone.ChangeZoneHandler
	observer   *observability.Observability
}

func NewChangeZoneHandler(
	cmdHandler *zone.ChangeZoneHandler,
	observer *observability.Observability,
) *ChangeZoneHandler {
	return &ChangeZoneHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//		@Summary  Change zone
//		@Tags    zones
//		@Accept    json
//		@Produce  json
//		@Param   params   path   GetZoneByIDRequest   true   "Search zone"
//		@Param     request  body   ChangeZoneRequest  true  "Zone data"
//		@Success  200        {object}  response.BaseResponseData[string]
//		@Failure      400  {object}  BaseResponse
//		@Failure      401  {object}  BaseResponse
//		@Failure      404  {object}  BaseResponse
//		@Failure      422  {object}  BaseResponse
//		@Router    /zones/{id} [patch]
//	 @Security		APIKeyAuth
func (h *ChangeZoneHandler) Handle(ctx *fiber.Ctx) error {
	var pathReq request.GetZoneByIDRequest
	if err := ctx.ParamsParser(&pathReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse id request")
		resp := response.NewErrorResponse("failed to parse id request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	var bodyReq request.ChangeZoneRequest
	if err := ctx.BodyParser(&bodyReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse body request")
		resp := response.NewErrorResponse("failed to parse body request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	var boundaries *[]app_model.ApplicationLocation
	if bodyReq.Boundaries != nil {
		boundaries = &[]app_model.ApplicationLocation{}
		for _, loc := range *bodyReq.Boundaries {
			*boundaries = append(*boundaries, app_model.ApplicationLocation{
				Latitude:  loc.Latitude,
				Longitude: loc.Longitude,
			})
		}
	}

	cmd := zone.ChangeZoneCommand{
		ID:         pathReq.ID,
		Name:       bodyReq.Name,
		Boundaries: boundaries,
		Safety:     bodyReq.Safety,
		UserID:     *identity.Sub,
	}

	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to change zone")
		return http.RespondError(ctx, err)
	}

	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
