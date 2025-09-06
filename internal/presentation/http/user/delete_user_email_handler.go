package user

//type DeleteUserEmailHandler struct {
//	qryHandler *user.DeleteUserEmailHandler
//	observer   *observability.Observability
//}
//
//func NewDeleteUserEmailHandler(
//	qryHandler *user.DeleteUserEmailHandler,
//	observer *observability.Observability,
//) *DeleteUserEmailHandler {
//	return &DeleteUserEmailHandler{
//		qryHandler: qryHandler,
//		observer:   observer,
//	}
//}
//

//func (h *DeleteUserEmailHandler) Handle(ctx *fiber.Ctx) error {
//	var req request.DeleteUserEmailRequest
//	err := ctx.ParamsParser(&req)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to parse request")
//		resp := response.NewErrorResponse("failed to parse request")
//		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(resp)
//	}
//
//	qry := user.DeleteUserEmailQuery{externalID: req.externalID}
//	err = h.qryHandler.Handle(ctx.UserContext(), qry)
//	if err != nil {
//		h.observer.Logger.Error().Err(err).Msg("failed to delete user email")
//		return response.RespondError(ctx, err)
//	}
//	resp := response.NewSuccessResponse("user email successfully deleted")
//	return ctx.Status(fiber.StatusOK).JSON(resp)
//}
