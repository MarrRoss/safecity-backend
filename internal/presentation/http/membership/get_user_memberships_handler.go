package membership

import (
	"awesomeProjectDDD/internal/application/membership"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetUserMembershipsHandler struct {
	qryHandler *membership.GetUserMembershipsHandler
	observer   *observability.Observability
}

func NewGetUserMembershipsHandler(
	qryHandler *membership.GetUserMembershipsHandler,
	observer *observability.Observability) *GetUserMembershipsHandler {
	return &GetUserMembershipsHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get user families with members
//	@Tags    families
//	@Accept    json
//	@Produce  json
//	@Success 200 {object} response.BaseResponseData[[]response.GetFamilyMembershipResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /families [get]
//	@Security		APIKeyAuth
func (h *GetUserMembershipsHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get identity")
		return http.RespondError(ctx, err)
	}

	qry := membership.GetUserMembershipsQuery{ID: *identity.Sub}
	memberships, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get user memberships")
		return http.RespondError(ctx, err)
	}

	membershipsResp := response.NewGetFamilyMembershipsResponse(memberships)
	resp := response.NewSuccessResponse[[]*response.GetFamilyMembershipResponse](membershipsResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
