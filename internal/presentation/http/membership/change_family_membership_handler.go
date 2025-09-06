package membership

import (
	"awesomeProjectDDD/internal/application/membership"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type ChangeFamilyMembershipHandler struct {
	cmdHandler *membership.ChangeFamilyMembershipHandler
	observer   *observability.Observability
}

func NewChangeFamilyMembershipHandler(
	cmdHandler *membership.ChangeFamilyMembershipHandler,
	observer *observability.Observability,
) *ChangeFamilyMembershipHandler {
	return &ChangeFamilyMembershipHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//		@Summary  Change membership role
//		@Tags    memberships
//		@Accept    json
//		@Produce  json
//	    @Param   params   path   GetFamilyMembershipByIDRequest   true   "Search membership"
//		@Param     request  body   ChangeFamilyMembershipRequest  true  "Membership data"
//		@Success  200        {object}  response.BaseResponseData[string]
//		@Failure      400  {object}  BaseResponse
//		@Failure      401  {object}  BaseResponse
//		@Failure      404  {object}  BaseResponse
//		@Failure      422  {object}  BaseResponse
//		@Router    /memberships/{id} [patch]
//		@Security		APIKeyAuth
func (h *ChangeFamilyMembershipHandler) Handle(ctx *fiber.Ctx) error {
	var pathReq request.GetFamilyMembershipByIDRequest
	if err := ctx.ParamsParser(&pathReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse id request")
		resp := response.NewErrorResponse("failed to parse id request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	var bodyReq request.ChangeFamilyMembershipRequest
	if err := ctx.BodyParser(&bodyReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse body request")
		resp := response.NewErrorResponse("failed to parse body request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := membership.ChangeFamilyMembershipCommand{
		ID:     pathReq.ID,
		RoleID: int(bodyReq.RoleID),
	}
	err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to change family membership")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
