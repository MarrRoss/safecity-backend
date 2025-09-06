package notification_log

import (
	"awesomeProjectDDD/internal/application/notification_log"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetNotificationLogsHandler struct {
	cmdHandler *notification_log.GetNotificationLogsHandler
}

func NewGetNotificationLogsHandler(
	cmdHandler *notification_log.GetNotificationLogsHandler) *GetNotificationLogsHandler {
	return &GetNotificationLogsHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *GetNotificationLogsHandler) Handle(ctx *fiber.Ctx) error {
	notifLogs, err := h.cmdHandler.Handle(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}

	resp := response.NewGetNotificationLogsResponse(notifLogs)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
