package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
)

type GetUserHandler struct {
	qryHandler *user.GetUserHandler
	observer   *observability.Observability
}

func NewGetUserHandler(qryHandler *user.GetUserHandler, observer *observability.Observability) *GetUserHandler {
	return &GetUserHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

func (h *GetUserHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.GetByUserIDRequest
	//err := ctx.ParamsParser(&req)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
	//	resp := response.NewErrorResponse("failed to parse request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//
	//qry := user.GetUserQuery{ID: req.ID.String()}
	//u, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get user")
	//	return http.RespondError(ctx, err)
	//}
	//
	//userResp := response.NewGetUserResponse(u)
	//resp := response.NewSuccessResponse(userResp)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
