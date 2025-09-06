package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetNotificationsSendersByReceiverHandler struct {
	cmdHandler *notification_setting.GetNotificationsSendersByReceiverHandler
}

func NewGetNotificationsSendersByReceiverHandler(
	cmdHandler *notification_setting.GetNotificationsSendersByReceiverHandler) *GetNotificationsSendersByReceiverHandler {
	return &GetNotificationsSendersByReceiverHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *GetNotificationsSendersByReceiverHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetNotificationsSendersByReceiverRequest
	req.UserID = ctx.Params("user_id")
	if req.UserID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "missing user_id"})
	}

	cmd := notification_setting.GetNotificationsSendersByReceiverCommand{
		ID: req.UserID,
	}
	notifySenders, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}

	resp := response.NewGetNotificationsSendersByReceiverResponse(notifySenders)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
