package user

import (
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetAllUsersHandler struct {
	qryHandler *user.GetAllUsersHandler
	observer   *observability.Observability
}

func NewGetAllUsersHandler(
	qryHandler *user.GetAllUsersHandler,
	observer *observability.Observability,
) *GetAllUsersHandler {
	return &GetAllUsersHandler{
		qryHandler: qryHandler,
		observer:   observer,
	}
}

// Handle
//
//	@Summary  Get users
//	@Tags    users
//	@Accept    json
//	@Produce  json
//	@Param   params query   GetUsersRequest   true   "Search users"
//	@Success 200 {object} response.BaseResponseData[[]response.GetUserResponse]
//	@Failure      400  {object}  BaseResponse
//	@Failure      401  {object}  BaseResponse
//	@Failure      404  {object}  BaseResponse
//	@Failure      422  {object}  BaseResponse
//	@Router    /users [get]
//	@Security		APIKeyAuth
func (h *GetAllUsersHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetUsersRequest
	err := ctx.QueryParser(&req)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
		resp := response.NewErrorResponse("failed to parse request")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}
	//h.observer.Logger.Info().Interface("req", req).Msg("got request")

	qry := user.GetAllUsersQuery{Login: req.Login}
	users, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get users")
		return http.RespondError(ctx, err)
	}

	userResp := response.NewGetUsersResponse(users)
	resp := response.NewSuccessResponse(userResp)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
