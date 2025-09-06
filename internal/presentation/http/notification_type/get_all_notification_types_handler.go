package notification_type

import (
	"awesomeProjectDDD/internal/application/notification_type"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAllNotificationTypesHandler struct {
	qryHandler *notification_type.GetAllNotificationTypesHandler
	observer   *observability.Observability
}

func NewGetAllNotificationTypesHandler(
	qryHandler *notification_type.GetAllNotificationTypesHandler,
	observer *observability.Observability,
) *GetAllNotificationTypesHandler {
	return &GetAllNotificationTypesHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

func (h *GetAllNotificationTypesHandler) Handle(ctx *fiber.Ctx) error {
	notificationTypes, err := h.qryHandler.Handle(ctx.UserContext())
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get notification types")
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}

	resp := response.NewGetNotificationTypesResponse(notificationTypes)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
