package membership

import (
	"awesomeProjectDDD/internal/application/membership"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddUserMembershipHandler struct {
	cmdHandler *membership.AddUserMembershipHandler
	observer   *observability.Observability
}

func NewAddUserMembershipHandler(
	cmdHandler *membership.AddUserMembershipHandler,
	observer *observability.Observability) *AddUserMembershipHandler {
	return &AddUserMembershipHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//		@Summary  Add membership from invitation
//		@Tags    invitations
//		@Accept    json
//		@Produce  json
//	    @Param   params   path   GetInvitationByIDRequest   true   "Search invitation"
//		@Success  200        {object}  response.BaseResponseData[response.AddFamilyMembershipResponse]
//		@Failure      400  {object}  BaseResponse
//	    @Failure      401  {object}  BaseResponse
//		@Failure      404  {object}  BaseResponse
//		@Failure      422  {object}  BaseResponse
//		@Router    /invitations/{id}/join [post]
//	    @Security		APIKeyAuth
func (h *AddUserMembershipHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	var pathReq request.GetInvitationByIDRequest
	if err := ctx.ParamsParser(&pathReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse id request")
		resp := response.NewErrorResponse("failed to parse id request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := membership.AddUserMembershipCommand{
		ID:     pathReq.ID,
		UserID: *identity.Sub,
	}
	id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add user membership from invitation")
		return http.RespondError(ctx, err)
	}

	idResp := response.NewAddFamilyMembershipResponse(id)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
