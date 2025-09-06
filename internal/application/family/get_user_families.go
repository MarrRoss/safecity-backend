package family

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type GetUserFamiliesQuery struct {
	ID string
}

type GetUserFamiliesHandler struct {
	familyStorage db.FamilyRepository
	userStorage   db.UserRepository
}

func NewGetUserFamiliesHandler(familyStorage db.FamilyRepository,
	userStorage db.UserRepository) *GetUserFamiliesHandler {
	return &GetUserFamiliesHandler{
		familyStorage: familyStorage,
		userStorage:   userStorage,
	}
}

func (h *GetUserFamiliesHandler) Handle(ctx context.Context,
	query GetUserFamiliesQuery) ([]*app_model.ApplicationFamily, error) {
	//userID, err := value_object.NewIDFromString(query.externalID)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid user id: %w", err)
	//}
	//user, err := h.userStorage.GetUser(ctx, userID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get user: %w", err)
	//}
	//families, err := h.familyStorage.GetFamiliesByAuthorID(ctx, user.externalID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get user families: %w", err)
	//}
	//if len(families) == 0 {
	//	return nil, errors.New("user doesn't have families")
	//}
	return app_model.NewApplicationFamilies(nil), nil
}
