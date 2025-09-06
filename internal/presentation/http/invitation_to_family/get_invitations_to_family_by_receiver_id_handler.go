package invitation_to_family

import (
	"awesomeProjectDDD/internal/application/invitation_to_family"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetInvitationsToFamilyByReceiverIdHandler struct {
	qryHandler *invitation_to_family.GetInvitationsToFamilyByReceiverIdHandler
	observer   *observability.Observability
}

func NewGetInvitationsToFamilyByReceiverIdHandler(
	qryHandler *invitation_to_family.GetInvitationsToFamilyByReceiverIdHandler,
	observer *observability.Observability,
) *GetInvitationsToFamilyByReceiverIdHandler {
	return &GetInvitationsToFamilyByReceiverIdHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get family invitations
//	@Tags    invitations
//	@Accept    json
//	@Produce  json
//	@Param     params      path   GetInvitationByIDRequest      true  "Search invitation"
//	@Success 200 {object} response.BaseResponseData[response.GetInvitationResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Router    /invitations/{id} [get]
func (h *GetInvitationsToFamilyByReceiverIdHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetInvitationByIDRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	qry := invitation_to_family.GetInvitationsToFamilyByReceiverIdQuery{
		InvitationID: req.ID,
	}
	invitations, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get user invitations")
		return http.RespondError(ctx, err)
	}

	invitationsResp := response.NewGetInvitationByReceiverResponse(invitations)
	resp := response.NewSuccessResponse[*response.GetInvitationResponse](invitationsResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
