package user

//type ChangeUserPasswordHandler struct {
//	cmdHandler *user.ChangeUserPasswordHandler
//	observer   *observability.Observability
//}
//
//func NewChangeUserPasswordHandler(
//	cmdHandler *user.ChangeUserPasswordHandler,
//	observer *observability.Observability,
//) *ChangeUserPasswordHandler {
//	return &ChangeUserPasswordHandler{
//		cmdHandler: cmdHandler,
//		observer:   observer,
//	}
//}
//
//func (h *ChangeUserPasswordHandler) Handle(ctx *fiber.Ctx) error {
//	var req request.ChangeUserPasswordRequest
//	err := ctx.BodyParser(&req)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
//		resp := response.NewErrorResponse("failed to parse request")
//		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
//	}
//	cmd := user.ChangeUserPasswordCommand{
//		externalID:       req.externalID,
//		Password: req.Password,
//	}
//	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
//	if err != nil {
//		h.observer.Logger.Error().Err(err).Msg("failed to change user password")
//		return response.RespondError(ctx, err)
//	}
//	resp := response.NewSuccessResponse[any](nil)
//	return ctx.Status(fiber.StatusOK).JSON(resp)
//}
