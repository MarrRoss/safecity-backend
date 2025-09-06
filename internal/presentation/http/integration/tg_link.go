package integration

import (
	"awesomeProjectDDD/internal/application/integration"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddTgLinkHandler struct {
	cmdHandler *integration.AddTgLinkHandler
	observer   *observability.Observability
}

func NewAddTgLinkHandler(
	cmdHandler *integration.AddTgLinkHandler,
	observer *observability.Observability) *AddTgLinkHandler {
	return &AddTgLinkHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add integration link
//	@Tags    integrations
//	@Accept    json
//	@Produce  json
//	@Success  200        {object}  response.BaseResponseData[string]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Router    /tg_link [post]
//	@Security		APIKeyAuth
func (h *AddTgLinkHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}
	cmd := integration.AddTgLinkCommand{
		UserID: *identity.Sub,
	}
	link, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add tg link")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse(link)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
