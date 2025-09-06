package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAvailableZonesHandler struct {
	qryHandler *notification_setting.GetAvailableZonesHandler
}

func NewGetAvailableZonesHandler(
	qryHandler *notification_setting.GetAvailableZonesHandler) *GetAvailableZonesHandler {
	return &GetAvailableZonesHandler{qryHandler: qryHandler}
}

func (h *GetAvailableZonesHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetAvailableZonesRequest
	req.ReceiverID = ctx.Params("receiver_id")
	if req.ReceiverID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "missing receiver_id"})
	}
	req.SenderID = ctx.Params("sender_id")
	if req.SenderID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "missing sender_id"})
	}

	qry := notification_setting.GetAvailableZonesQuery{
		ReceiverID: req.ReceiverID,
		SenderID:   req.SenderID,
	}
	zones, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	resp := response.NewGetZonesResponse(zones)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
