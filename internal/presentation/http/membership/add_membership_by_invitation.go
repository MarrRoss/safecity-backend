package membership

import (
	"awesomeProjectDDD/internal/application/membership"
	"awesomeProjectDDD/internal/presentation/http/request"
	"github.com/gofiber/fiber/v2"
)

type AddMembershipByInvitationHandler struct {
	cmdHandler *membership.AddMembershipByInvitationHandler
}

func NewAddMembershipByInvitationHandler(
	cmdHandler *membership.AddMembershipByInvitationHandler) *AddMembershipByInvitationHandler {
	return &AddMembershipByInvitationHandler{
		cmdHandler: cmdHandler,
	}
}

func (h *AddMembershipByInvitationHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddMembershipByInvitationRequest
	req.ID = ctx.Params("invitation_id")
	if req.ID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "missing invitation_id"})
	}

	cmd := membership.AddMembershipByInvitationCommand{
		ID: req.ID,
	}
	id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).
		JSON(map[string]string{"id": id})
}
