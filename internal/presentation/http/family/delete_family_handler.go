package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/presentation/http/request"
	"github.com/gofiber/fiber/v2"
)

type DeleteFamilyHandler struct {
	qryHandler *family.DeleteFamilyHandler
}

func NewDeleteFamilyHandler(qryHandler *family.DeleteFamilyHandler) *DeleteFamilyHandler {
	return &DeleteFamilyHandler{
		qryHandler: qryHandler,
	}
}

func (dfh *DeleteFamilyHandler) Handle(ctx *fiber.Ctx) error {
	var req request.DeleteFamilyRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "enable to parse request"})
	}

	qry := family.DeleteFamilyQuery{ID: req.ID}
	err = dfh.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	return ctx.
		Status(fiber.StatusOK).
		JSON(map[string]string{"status": "ok", "detail": "family successfully deleted"})
}
