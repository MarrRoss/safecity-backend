package membership

import (
	"awesomeProjectDDD/internal/application/membership"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type DeleteUserMembershipHandler struct {
	cmdHandler *membership.DeleteUserMembershipHandler
	observer   *observability.Observability
}

func NewDeleteUserMembershipHandler(
	cmdHandler *membership.DeleteUserMembershipHandler,
	observer *observability.Observability,
) *DeleteUserMembershipHandler {
	return &DeleteUserMembershipHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

func (h *DeleteUserMembershipHandler) Handle(ctx *fiber.Ctx) error {
	var req request.DeleteUserMembershipRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	cmd := membership.DeleteUserMembershipCommand{
		ID:     req.ID,
		UserID: *identity.Sub,
	}
	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to delete user membership")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
