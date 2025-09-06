package messenger_type

//type GetAllMessengerTypesHandler struct {
//	qryHandler *messenger_type.GetAllMessengerTypesHandler
//	observer   *observability.Observability
//}
//
//func NewGetAllMessengerTypesHandler(
//	qryHandler *messenger_type.GetAllMessengerTypesHandler,
//	observer *observability.Observability,
//) *GetAllMessengerTypesHandler {
//	return &GetAllMessengerTypesHandler{
//		qryHandler: qryHandler,
//		observer:   observer,
//	}
//}
//
//func (h *GetAllMessengerTypesHandler) Handle(ctx *fiber.Ctx) error {
//	messengerTypes, err := h.qryHandler.Handle(ctx.UserContext())
//	if err != nil {
//		h.observer.Logger.Error().Err(err).Msg("failed to get messenger types")
//		return ctx.Status(fiber.StatusInternalServerError).
//			JSON(map[string]string{"status": "error", "detail": err.Error()})
//	}
//
//	resp := response.NewGetMessengerTypesResponse(messengerTypes)
//	return ctx.Status(fiber.StatusOK).JSON(resp)
//}
