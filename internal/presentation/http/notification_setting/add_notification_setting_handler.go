package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddNotificationSettingHandler struct {
	cmdHandler *notification_setting.AddNotificationSettingHandler
	observer   *observability.Observability
}

func NewAddNotificationSettingHandler(
	cmdHandler *notification_setting.AddNotificationSettingHandler,
	observer *observability.Observability,
) *AddNotificationSettingHandler {
	return &AddNotificationSettingHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add notification setting
//	@Tags    notifications
//	@Accept    json
//	@Produce  json
//	@Param     request  body   AddNotificationSettingRequest  true  "Notification setting data"
//	@Success  200        {object}  response.BaseResponseData[response.AddNotificationSettingResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /notifications_settings [post]
//	@Security		APIKeyAuth
func (h *AddNotificationSettingHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddNotificationSettingRequest
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}
	err = ctx.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := notification_setting.AddNotificationSettingCommand{
		FrequencyID: req.FrequencyID,
		EventType:   string(req.EventType),
		ReceiverID:  *identity.Sub,
		SenderID:    req.SenderID,
		ZoneID:      req.ZoneID,
		Battery:     req.Battery,
	}
	id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add notification setting")
		return http.RespondError(ctx, err)
	}
	idResp := response.NewAddNotificationSettingResponse(id)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
