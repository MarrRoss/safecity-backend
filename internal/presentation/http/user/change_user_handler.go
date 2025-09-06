package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
)

type ChangeUserHandler struct {
	cmdHandler *user.ChangeUserHandler
	observer   *observability.Observability
}

func NewChangeUserHandler(
	cmdHandler *user.ChangeUserHandler,
	observer *observability.Observability,
) *ChangeUserHandler {
	return &ChangeUserHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

func (h *ChangeUserHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.ChangeUserRequest
	//if err := ctx.ParamsParser(&req); err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse id request")
	//	resp := response.NewErrorResponse("failed to parse id request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//if err := ctx.BodyParser(&req); err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
	//	resp := response.NewErrorResponse("failed to parse request")
	//	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	//}
	//
	//cmd := user.ChangeUserCommand{
	//	ID:        req.ID,
	//	FirstName: req.FirstName,
	//	LastName:  req.LastName,
	//}
	//err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to change user")
	//	return http.RespondError(ctx, err)
	//}
	//resp := response.NewSuccessResponse[any](nil)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
