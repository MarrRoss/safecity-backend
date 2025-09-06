package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetMeHandler struct {
	qryHandler *user.GetUserHandler
	observer   *observability.Observability
}

func NewGetMeHandler(qryHandler *user.GetUserHandler, observer *observability.Observability) *GetMeHandler {
	return &GetMeHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get me
//	@Tags    users
//	@Accept    json
//	@Produce  json
//	@Success  200        {object}  response.BaseResponseData[response.GetUserResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Router    /me [get]
//	@Security		APIKeyAuth
func (h *GetMeHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	qry := user.GetUserQuery{ID: *identity.Sub}
	u, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get user")
		return http.RespondError(ctx, err)
	}

	userResp := response.NewGetUserResponse(u)
	resp := response.NewSuccessResponse(userResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
