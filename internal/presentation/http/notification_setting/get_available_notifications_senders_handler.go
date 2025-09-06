package notification_setting

import (
	"awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAvailableNotificationSendersHandler struct {
	qryHandler *notification_setting.GetAvailableNotificationSendersHandler
	observer   *observability.Observability
}

func NewGetAvailableNotificationSendersHandler(
	qryHandler *notification_setting.GetAvailableNotificationSendersHandler,
	observer *observability.Observability,
) *GetAvailableNotificationSendersHandler {
	return &GetAvailableNotificationSendersHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get available notification senders
//	@Tags    notifications
//	@Accept    json
//	@Produce  json
//	@Param     params      path   GetAvailableNotificationSendersRequest      true  "Search available notification senders"
//	@Success  200        {object}  response.BaseResponseData[[]response.GetUserResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /families/{id}/senders [get]
//	@Security		APIKeyAuth
func (h *GetAvailableNotificationSendersHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	var req request.GetAvailableNotificationSendersRequest
	err = ctx.ParamsParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	qry := notification_setting.GetAvailableNotificationSendersQuery{
		ReceiverID: *identity.Sub,
		FamilyID:   req.FamilyID,
	}
	senders, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get available notification senders")
		return http.RespondError(ctx, err)
	}

	sendersResp := response.NewGetUsersResponse(senders)
	resp := response.NewSuccessResponse[[]*response.GetUserResponse](sendersResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
