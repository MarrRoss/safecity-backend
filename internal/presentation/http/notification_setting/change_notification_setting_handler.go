package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http/request"
	"github.com/gofiber/fiber/v2"
)

type ChangeNotificationSettingHandler struct {
	cmdHandler *notification_setting.ChangeNotificationSettingHandler
	observer   *observability.Observability
}

func NewChangeNotificationSettingHandler(
	cmdHandler *notification_setting.ChangeNotificationSettingHandler,
	observer *observability.Observability,
) *ChangeNotificationSettingHandler {
	return &ChangeNotificationSettingHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

func (h *ChangeNotificationSettingHandler) Handle(ctx *fiber.Ctx) error {
	var req request.ChangeNotificationSettingRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "error parsing request body"})
	}

	cmd := notification_setting.ChangeNotificationSettingCommand{
		ID:             req.ID,
		FrequencyID:    req.FrequencyID,
		NotifyTypesIDs: req.NotTypesIDs,
		MesTypesIDs:    req.MesTypesIDs,
	}
	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to change notification setting")
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).
		JSON(map[string]string{"status": "ok", "detail": "notification setting successfully changed"})
}
