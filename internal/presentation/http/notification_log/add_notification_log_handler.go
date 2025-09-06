package notification_log

import (
	"awesomeProjectDDD/internal/application/notification_log"
	"awesomeProjectDDD/internal/presentation/http/request"
	"github.com/gofiber/fiber/v2"
)

type AddNotificationLogHandler struct {
	cmdHandler *notification_log.AddNotificationLogHandler
}

func NewAddNotificationLogHandler(
	cmdHandler *notification_log.AddNotificationLogHandler) *AddNotificationLogHandler {
	return &AddNotificationLogHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *AddNotificationLogHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddNotificationLogRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "error parsing request body"})
	}

	cmd := notification_log.AddNotificationLogCommand{
		NotificationSettingID: req.NotificationSettingID,
		SendTime:              req.SendTime,
		Context:               req.Context,
	}
	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).
		JSON(map[string]string{"status": "notification log successfully created", "detail": ""})

}
