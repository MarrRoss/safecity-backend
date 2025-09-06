package user_location

import (
	"awesomeProjectDDD/internal/adapter/telegram"
	"awesomeProjectDDD/internal/application/user_location"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddUserLocationHandler struct {
	cmdHandler      *user_location.AddUserLocationHandler
	observer        *observability.Observability
	telegramService *telegram.Service
}

func NewAddUserLocationHandler(
	cmdHandler *user_location.AddUserLocationHandler,
	observer *observability.Observability,
	telegramService *telegram.Service,
) *AddUserLocationHandler {
	return &AddUserLocationHandler{
		cmdHandler:      cmdHandler,
		observer:        observer,
		telegramService: telegramService,
	}
}

// Handle
//
//	@Summary  Add user location
//	@Tags    locations
//	@Accept    json
//	@Produce  json
//	@Param     request  body   AddUserLocationRequest  true  "User Location data"
//	@Success  200        {object}  response.BaseResponseData[response.AddUserLocationResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /users/locations [post]
//	@Security		APIKeyAuth
func (h *AddUserLocationHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}
	var req request.AddUserLocationRequest
	err = ctx.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := user_location.AddUserLocationCommand{
		UserID:    *identity.Sub,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Battery:   req.Battery,
	}
	locationLogID, notifyLogsIDs, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add user location")
		return http.RespondError(ctx, err)
	}

	idResp := response.NewAddUserLocationResponse(locationLogID, notifyLogsIDs)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
