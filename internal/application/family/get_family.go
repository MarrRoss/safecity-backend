package family

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type GetFamilyQuery struct {
	ID string
}

type GetFamilyHandler struct {
	storage db.FamilyRepository
}

func NewGetFamilyHandler(storage db.FamilyRepository) *GetFamilyHandler {
	return &GetFamilyHandler{
		storage: storage,
	}
}

func (h *GetFamilyHandler) Handle(ctx context.Context, query GetFamilyQuery) (*app_model.ApplicationFamily, error) {
	id, err := value_object.NewIDFromString(query.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid family id: %w", err)
	}
	family, err := h.storage.GetFamily(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get family: %w", err)
	}

	return app_model.NewApplicationFamily(family), nil
}
