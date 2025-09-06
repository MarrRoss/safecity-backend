package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http/request"
	"github.com/gofiber/fiber/v2"
)

type DeleteNotificationSettingHandler struct {
	cmdHandler *notification_setting.DeleteNotificationSettingHandler
	observer   *observability.Observability
}

func NewDeleteNotificationSettingHandler(
	cmdHandler *notification_setting.DeleteNotificationSettingHandler,
	observer *observability.Observability,
) *DeleteNotificationSettingHandler {
	return &DeleteNotificationSettingHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

func (h *DeleteNotificationSettingHandler) Handle(ctx *fiber.Ctx) error {
	var req request.DeleteNotificationSettingRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "error parsing request"})
	}

	cmd := notification_setting.DeleteNotificationSettingCommand{
		ID: req.ID,
	}
	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to delete notification setting")
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).
		JSON(map[string]string{"status": "ok", "detail": "notification setting successfully deleted"})
}
