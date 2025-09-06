package role

import (
	"awesomeProjectDDD/internal/application/role"
	"awesomeProjectDDD/internal/presentation/http"
	"awesomeProjectDDD/internal/presentation/http/request"
	"awesomeProjectDDD/internal/presentation/http/response"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type GetRoleHandler struct {
	qryHandler *role.GetRoleHandler
}

func NewGetRoleHandler(qryHandler *role.GetRoleHandler) *GetRoleHandler {
	return &GetRoleHandler{
		qryHandler: qryHandler,
	}
}

func (h *GetRoleHandler) Handle(ctx *fiber.Ctx) error {
	identity, err := http.GetIdentity(ctx)
	if err != nil {
		return http.RespondError(ctx, err)
	}
	fmt.Printf("identity: %v\n", identity)

	var req request.GetRoleByIDRequest
	err = ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(map[string]string{"status": "error", "detail": "error parsing request"})
	}
	qry := role.GetRoleQuery{ID: req.ID, ExternalUserID: *identity.Sub}

	r, err := h.qryHandler.Handle(ctx.UserContext(), qry)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(map[string]string{"status": "error", "detail": err.Error()})
	}
	resp := response.NewGetRoleByIDResponse(r)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
