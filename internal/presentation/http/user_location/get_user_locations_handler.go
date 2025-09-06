package user_location

import (
	"awesomeProjectDDD/internal/application/user_location"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetUserLocationsHandler struct {
	qryHandler *user_location.GetUserLocationsHandler
	observer   *observability.Observability
}

func NewGetUserLocationsHandler(
	qryHandler *user_location.GetUserLocationsHandler,
	observer *observability.Observability,
) *GetUserLocationsHandler {
	return &GetUserLocationsHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get users locations
//	@Tags    locations
//	@Accept    json
//	@Produce  json
//	@Success  200        {object}  response.BaseResponseData[[]response.GetUserLocationResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Router    /users/locations [get]
//	@Security		APIKeyAuth
func (h *GetUserLocationsHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	qry := user_location.GetUserLocationsQuery{UserID: *identity.Sub}
	locations, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get users locations")
		return http.RespondError(ctx, err)
	}

	locationsResp := response.NewGetUserLocationResponses(locations)
	resp := response.NewSuccessResponse[[]*response.GetUserLocationResponse](locationsResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
