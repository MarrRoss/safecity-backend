package membership

import (
	"awesomeProjectDDD/internal/application/membership"
	"github.com/gofiber/fiber/v2"
)

type GetMembershipsByFamilyHandler struct {
	qryHandler *membership.GetMembershipsByFamilyHandler
}

func NewGetMembershipsByFamilyHandler(
	qryHandler *membership.GetMembershipsByFamilyHandler) *GetMembershipsByFamilyHandler {
	return &GetMembershipsByFamilyHandler{
		qryHandler: qryHandler,
	}
}

func (h *GetMembershipsByFamilyHandler) Handle(ctx *fiber.Ctx) error {
	//var req request.GetMembershipsByFamilyRequest
	//req.externalID = ctx.Params("family_id")
	//if req.externalID == "" {
	//	return ctx.Status(fiber.StatusBadRequest).
	//		JSON(map[string]string{"status": "error", "detail": "missing family_id"})
	//}
	//
	//qry := membership.GetMembershipsByFamilyQuery{externalID: req.externalID}
	//familyMemberships, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	//if err != nil {
	//	return ctx.Status(fiber.StatusInternalServerError).
	//		JSON(map[string]string{"status": "error", "detail": err.Error()})
	//}
	//
	//resp := response.NewGetMembershipsResponse(familyMemberships)
	//return ctx.Status(fiber.StatusOK).JSON(resp)
	return nil
}
