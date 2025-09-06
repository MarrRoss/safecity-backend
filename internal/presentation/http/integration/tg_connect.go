package integration

import (
	"awesomeProjectDDD/internal/application/integration"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddTgConnectHandler struct {
	cmdHandler *integration.AddTgConnectHandler
	observer   *observability.Observability
}

func NewAddTgConnectHandler(
	cmdHandler *integration.AddTgConnectHandler,
	observer *observability.Observability) *AddTgConnectHandler {
	return &AddTgConnectHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add integration connect
//	@Tags    integrations
//	@Accept    json
//	@Produce  json
//	@Param     request  body   AddTgConnectRequest  true  "Integration connect data"
//	@Success  200        {object}  response.BaseResponseData[string]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /tg_connect [post]
func (h *AddTgConnectHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddTgConnectRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	cmd := integration.AddTgConnectCommand{
		UserID:     req.UserID,
		TelegramID: req.TgID,
	}
	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add tg link")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
