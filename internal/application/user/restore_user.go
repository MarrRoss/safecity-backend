package user

import (
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type RestoreUserQuery struct {
	ID string
}

type RestoreUserHandler struct {
	userStore db.UserRepository
}

func NewRestoreUserHandler(storage db.UserRepository) *RestoreUserHandler {
	return &RestoreUserHandler{
		userStore: storage,
	}
}

func (h *RestoreUserHandler) Handle(ctx context.Context, query RestoreUserQuery) error {
	id, err := value_object.NewIDFromString(query.ID)
	if err != nil {
		return fmt.Errorf("invalid user id: %s", query.ID)
	}
	user, err := h.userStore.GetUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	//err = user.RestoreExistence()
	//if err != nil {
	//	return fmt.Errorf("failed to restore user existence: %w", err)
	//}
	err = h.userStore.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
