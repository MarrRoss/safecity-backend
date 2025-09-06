package user

//type ChangeUserPhoneHandler struct {
//	cmdHandler *user.ChangeUserPhoneHandler
//	observer   *observability.Observability
//}
//
//func NewChangeUserPhoneHandler(
//	cmdHandler *user.ChangeUserPhoneHandler,
//	observer *observability.Observability,
//) *ChangeUserPhoneHandler {
//	return &ChangeUserPhoneHandler{
//		cmdHandler: cmdHandler,
//		observer:   observer,
//	}
//}
//

//func (h *ChangeUserPhoneHandler) Handle(ctx *fiber.Ctx) error {
//	var req request.ChangeUserPhoneRequest
//	err := ctx.BodyParser(&req)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
//		resp := response.NewErrorResponse("failed to parse request")
//		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
//	}
//	cmd := user.ChangeUserPhoneCommand{
//		externalID:    req.externalID,
//		Phone: req.Phone,
//	}
//	err = h.cmdHandler.Handle(ctx.UserContext(), cmd)
//	if err != nil {
//		h.observer.Logger.Error().Err(err).Msg("failed to change user phone")
//		return response.RespondError(ctx, err)
//	}
//	resp := response.NewSuccessResponse("user phone successfully changed")
//	return ctx.Status(fiber.StatusOK).JSON(resp)
//}
