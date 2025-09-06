package invitation_to_family

import (
	"awesomeProjectDDD/internal/application/invitation_to_family"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddInvitationToFamilyHandler struct {
	cmdHandler *invitation_to_family.AddInvitationToFamilyHandler
	observer   *observability.Observability
}

func NewAddInvitationToFamilyHandler(
	cmdHandler *invitation_to_family.AddInvitationToFamilyHandler,
	observer *observability.Observability,
) *AddInvitationToFamilyHandler {
	return &AddInvitationToFamilyHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add invitation to family
//	@Tags    invitations
//	@Accept    json
//	@Produce  json
//	@Param     request  body   AddInvitationToFamilyRequest  true  "Invitation data"
//	@Success  200        {object}  response.BaseResponseData[response.AddInvitationResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /invitations [post]
//	@Security		APIKeyAuth
func (h *AddInvitationToFamilyHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddInvitationToFamilyRequest
	err := ctx.BodyParser(&req)
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

	cmd := invitation_to_family.AddInvitationToFamilyCommand{
		AuthorID: *identity.Sub,
		RoleID:   req.RoleID,
		FamilyID: req.FamilyID,
	}
	id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add invitation to family")
		return http.RespondError(ctx, err)
	}
	idResp := response.NewAddInvitationResponse(id)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
