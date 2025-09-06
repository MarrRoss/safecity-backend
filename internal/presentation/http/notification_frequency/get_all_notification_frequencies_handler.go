package notification_frequency

import (
	"awesomeProjectDDD/internal/application/notification_frequency"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAllNotificationFrequenciesHandler struct {
	qryHandler *notification_frequency.GetAllNotificationFrequenciesHandler
	observer   *observability.Observability
}

func NewGetAllNotificationFrequenciesHandler(
	qryHandler *notification_frequency.GetAllNotificationFrequenciesHandler,
	observer *observability.Observability,
) *GetAllNotificationFrequenciesHandler {
	return &GetAllNotificationFrequenciesHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get frequencies
//	@Tags    notifications
//	@Accept    json
//	@Produce  json
//	@Success 200 {object} response.BaseResponseData[[]response.GetNotificationFrequencyResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Router    /frequencies [get]
//	@Security		APIKeyAuth
func (h *GetAllNotificationFrequenciesHandler) Handle(ctx *fiber.Ctx) error {
	frequencies, err := h.qryHandler.Handle(ctx.UserContext())
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get notification frequencies")
		return http.RespondError(ctx, err)
	}

	frequenciesResp := response.NewGetNotificationFrequenciesResponse(frequencies)
	resp := response.NewSuccessResponse[[]*response.GetNotificationFrequencyResponse](frequenciesResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
