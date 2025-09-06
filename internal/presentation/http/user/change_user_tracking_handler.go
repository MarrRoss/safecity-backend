package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type ChangeUserTrackingHandler struct {
	cmdHandler *user.ChangeUserTrackingHandler
	observer   *observability.Observability
}

func NewChangeUserTrackingHandler(
	cmdHandler *user.ChangeUserTrackingHandler,
	observer *observability.Observability,
) *ChangeUserTrackingHandler {
	return &ChangeUserTrackingHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Change user tracking
//	@Tags    users
//	@Accept    json
//	@Produce  json
//	@Success  200        {object}  response.BaseResponseData[string]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Router    /users/tracking [put]
//	@Security		APIKeyAuth
func (h *ChangeUserTrackingHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	cmd := user.ChangeUserTrackingCommand{
		ID: *identity.Sub,
	}
	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to change user tracking")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
