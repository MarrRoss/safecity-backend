package user

import (
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"github.com/google/uuid"
)

type DeleteUserQuery struct {
	ID uuid.UUID
}

type DeleteUserHandler struct {
	usersStore       db.UserRepository
	membershipsStore db.FamilyMembershipRepository
	observer         *observability.Observability
}

func NewDeleteUserHandler(
	usersStorage db.UserRepository,
	membershipsStore db.FamilyMembershipRepository,
	observer *observability.Observability,
) *DeleteUserHandler {
	return &DeleteUserHandler{
		usersStore:       usersStorage,
		membershipsStore: membershipsStore,
		observer:         observer,
	}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, query DeleteUserQuery) error {
	//id, err := value_object.NewIDFromString(query.ID.String())
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid user id")
	//	return err
	//}
	//user, err := h.usersStore.GetUser(ctx, id)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrUserNotFound) {
	//		h.observer.Logger.Trace().Err(err).Msg("user not found")
	//		return err
	//	}
	//
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting user")
	//		return err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
	//	return err
	//}
	//err = user.StopExistence()
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("failed to stop user existence")
	//	return err
	//}
	//
	////memberships, err := h.membershipsStore.GetMembershipsByUserID(id)
	////if err != nil {
	////	return fmt.Errorf("failed to get memberships: %w", err)
	////}
	////for _, membership := range memberships {
	////	membershipID, err := value_object.NewIDFromString(membership.ID)
	////	if err != nil {
	////		return fmt.Errorf("invalid membership id: %s", membership.ID)
	////	}
	////	err = h.membershipsStore.DeleteMembership(membershipID)
	////	if err != nil {
	////		return fmt.Errorf("failed to delete membership: %w", err)
	////	}
	////}
	//
	//err = h.usersStore.UpdateUser(ctx, user)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while updating user")
	//		return err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while updating user")
	//	return err
	//}
	return nil
}
