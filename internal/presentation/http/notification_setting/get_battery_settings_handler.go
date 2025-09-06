package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetBatterySettingsHandler struct {
	cmdHandler *notification_setting.GetBatterySettingsHandler
	observer   *observability.Observability
}

func NewGetBatterySettingsHandler(
	cmdHandler *notification_setting.GetBatterySettingsHandler,
	observer *observability.Observability,
) *GetBatterySettingsHandler {
	return &GetBatterySettingsHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get user battery notify settings
//	@Tags    notifications
//	@Accept    json
//	@Produce  json
//	@Success  200        {object}  response.BaseResponseData[[]response.GetNotificationSettingResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /battery/notifications_settings [get]
//	@Security		APIKeyAuth
func (h *GetBatterySettingsHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	cmd := notification_setting.GetBatterySettingsQuery{
		UserID: *identity.Sub,
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
