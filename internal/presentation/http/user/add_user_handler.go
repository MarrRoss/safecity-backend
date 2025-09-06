package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type AddUserHandler struct {
	cmdHandler *user.AddUserHandler
	observer   *observability.Observability
}

func NewAddUserHandler(cmdHandler *user.AddUserHandler, observer *observability.Observability) *AddUserHandler {
	return &AddUserHandler{
		cmdHandler: cmdHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Add user
//	@Tags    users
//	@Accept    json
//	@Produce  json
//	@Param     request  body   AddUserRequest  true  "User data"
//	@Success  200        {object}  response.BaseResponseData[response.AddUserResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /users [post]
func (h *AddUserHandler) Handle(ctx *fiber.Ctx) error {
	var req request.AddUserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	cmd := user.AddUserCommand{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Login:      req.Login,
		ExternalID: req.ExternalID,
	}
	id, err := h.cmdHandler.Handle(ctx.UserContext(), cmd)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to add user")
		return http.RespondError(ctx, err)
	}
	idResp := response.NewAddUserResponse(id)
	resp := response.NewSuccessResponse(idResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
