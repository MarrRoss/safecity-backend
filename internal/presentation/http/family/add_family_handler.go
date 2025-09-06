package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddFamilyHandler struct {
	cmdHandler *family.AddFamilyHandler
	observer   *observability.Observability
}

func NewAddFamilyHandler(
	cmdHandler *family.AddFamilyHandler,
	observer *observability.Observability) *AddFamilyHandler {
	return &AddFamilyHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add family
//	@Tags    families
//	@Accept    json
//	@Produce  json
//	@Param     request  body   AddFamilyRequest  true  "Family data"
//	@Success  200        {object}  response.BaseResponseData[response.AddFamilyResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /families [post]
//	@Security		APIKeyAuth
func (h *AddFamilyHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddFamilyRequest
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
	cmd := family.AddFamilyCommand{
		Name:   req.Name,
		UserID: *identity.Sub,
	}
	familyID, membershipID, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add family")
		return http.RespondError(ctx, err)
	}
	idResp := response.NewAddFamilyResponse(familyID, membershipID)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
