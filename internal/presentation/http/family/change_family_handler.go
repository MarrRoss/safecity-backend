package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type ChangeFamilyHandler struct {
	cmdHandler *family.ChangeFamilyHandler
	observer   *observability.Observability
}

func NewChangeFamilyHandler(
	cmdHandler *family.ChangeFamilyHandler,
	observer *observability.Observability,
) *ChangeFamilyHandler {
	return &ChangeFamilyHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//		@Summary  Change family
//		@Tags    families
//		@Accept    json
//		@Produce  json
//	    @Param   params   path   GetFamilyByIDRequest   true   "Search family"
//		@Param     request  body   ChangeFamilyRequest  true  "Family data"
//		@Success  200        {object}  response.BaseResponseData[string]
//		@Failure      400  {object}  BaseResponse
//		@Failure      401  {object}  BaseResponse
//		@Failure      404  {object}  BaseResponse
//		@Failure      422  {object}  BaseResponse
//		@Router    /families/{id} [patch]
//		@Security		APIKeyAuth
func (h *ChangeFamilyHandler) Handle(ctx *fiber.Ctx) error {
	var pathReq request.GetFamilyByIDRequest
	if err := ctx.ParamsParser(&pathReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse id request")
		resp := response.NewErrorResponse("failed to parse id request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	var bodyReq request.ChangeFamilyRequest
	if err := ctx.BodyParser(&bodyReq); err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse body request")
		resp := response.NewErrorResponse("failed to parse body request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := family.ChangeFamilyCommand{
		ID:   pathReq.ID,
		Name: bodyReq.Name,
	}
	err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to change family")
		return http.RespondError(ctx, err)
	}
	resp := response.NewSuccessResponse[any](nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
