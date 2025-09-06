package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
)

type DeleteUserHandler struct {
	qryHandler *user.DeleteUserHandler
	observer   *observability.Observability
}

func NewDeleteUserHandler(
	qryHandler *user.DeleteUserHandler,
	observer *observability.Observability,
) *DeleteUserHandler {
	return &DeleteUserHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

func (h *DeleteUserHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.DeleteUserRequest
	//err := ctx.ParamsParser(&req)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
	//	resp := response.NewErrorResponse("failed to parse request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//
	//qry := user.DeleteUserQuery{ID: req.ID}
	//err = h.qryHandler.Handle(ctx.UserContext(), qry)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to delete user")
	//	return http.RespondError(ctx, err)
	//}
	//resp := response.NewSuccessResponse[any](nil)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
