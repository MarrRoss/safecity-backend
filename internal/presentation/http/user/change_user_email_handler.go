package user

//type ChangeUserEmailHandler struct {
//	cmdHandler *user.ChangeUserEmailHandler
//	observer   *observability.Observability
//}
//
//func NewChangeUserEmailHandler(
//	cmdHandler *user.ChangeUserEmailHandler,
//	observer *observability.Observability,
//) *ChangeUserEmailHandler {
//	return &ChangeUserEmailHandler{
//		cmdHandler: cmdHandler,
//		observer:   observer,
//	}
//}
//

//func (h *ChangeUserEmailHandler) Handle(ctx *fiber.Ctx) error {
//	var req request.ChangeUserEmailRequest
//	err := ctx.BodyParser(&req)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
//		resp := response.NewErrorResponse("failed to parse request")
//		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
//	}
//	cmd := user.ChangeUserEmailCommand{
//		externalID:    req.externalID,
//		Email: req.Email,
//	}
//	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
//	if err != nil {
//		h.observer.Logger.Error().Err(err).Msg("failed to change user email")
//		return response.RespondError(ctx, err)
//	}
//	resp := response.NewSuccessResponse[any](nil)
//	return ctx.Status(fiber.StatusOK).JSON(resp)
//}
