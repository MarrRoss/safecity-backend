package family

import (
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type GetFamilyHandler struct {
	qryHandler *family.GetFamilyHandler
}

func NewGetFamilyHandler(qryHandler *family.GetFamilyHandler) *GetFamilyHandler {
	return &GetFamilyHandler{
		qryHandler: qryHandler,
	}
}

func (gh *GetFamilyHandler) Handle(ctx *fiber.Ctx) error {
	var req request.GetFamilyByIDRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "error parsing request"})
	}
	qry := family.GetFamilyQuery{ID: req.ID.String()}
	f, err := gh.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	resp := response.NewGetFamilyResponse(f)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
