package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetNotificationSettingsHandler struct {
	cmdHandler *notification_setting.GetNotificationSettingsHandler
	observer   *observability.Observability
}

func NewGetNotificationSettingsHandler(
	cmdHandler *notification_setting.GetNotificationSettingsHandler,
	observer *observability.Observability,
) *GetNotificationSettingsHandler {
	return &GetNotificationSettingsHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get user notify settings by zone
//	@Tags    notifications
//	@Accept    json
//	@Produce  json
//	@Param     params      path   GetNotificationSettingsRequest      true  "Search notify settings"
//	@Success  200        {object}  response.BaseResponseData[[]response.GetNotificationSettingResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /zones/{id}/notifications_settings [get]
//	@Security		APIKeyAuth
func (h *GetNotificationSettingsHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetNotificationSettingsRequest
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}
	err = ctx.ParamsParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := notification_setting.GetNotificationSettingsQuery{
		UserID: *identity.Sub,
		ZoneID: req.ID,
	}
	notifySettings, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get notification settings")
		return http.RespondError(ctx, err)
	}

	notifySettingsResp := response.NewGetNotificationSettingsResponse(notifySettings)
	resp := response.NewSuccessResponse[[]*response.GetNotificationSettingResponse](notifySettingsResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
