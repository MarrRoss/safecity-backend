package notification_log

import (
	"awesomeProjectDDD/internal/application/notification_log"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetUserNotificationLogsHandler struct {
	cmdHandler *notification_log.GetUserNotificationLogsHandler
}

func NewGetUserNotificationLogsHandler(
	cmdHandler *notification_log.GetUserNotificationLogsHandler) *GetUserNotificationLogsHandler {
	return &GetUserNotificationLogsHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *GetUserNotificationLogsHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetUserNotificationLogsRequest
	req.UserID = ctx.Params("user_id")
	if req.UserID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "missing user_id"})
	}

	qry := notification_log.GetUserNotificationLogsQuery{UserID: req.UserID}
	notifLogs, err := h.cmdHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}

	resp := response.NewGetNotificationLogsResponse(notifLogs)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
