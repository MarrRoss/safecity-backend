package user

//type DeleteUserEmailQuery struct {
//	ID uuid.UUID
//}
//
//type DeleteUserEmailHandler struct {
//	userStore db.UserRepository
//	observer  *observability.Observability
//}
//
//func NewDeleteUserEmailHandler(
//	storage db.UserRepository,
//	observer *observability.Observability,
//) *DeleteUserEmailHandler {
//	return &DeleteUserEmailHandler{
//		userStore: storage,
//		observer:  observer,
//	}
//}
//
//func (h *DeleteUserEmailHandler) Handle(ctx context.Context, query DeleteUserEmailQuery) error {
//	id, err := value_object.NewIDFromString(query.ID.String())
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("invalid user id")
//		return err
//	}
//	user, err := h.userStore.GetUser(ctx, id)
//	if err != nil {
//		if errors.Is(err, domain.ErrUserNotFound) {
//			h.observer.Logger.Trace().Err(err).Msg("user not found")
//			return err
//		}
//
//		if errors.Is(err, adapter.ErrStorage) {
//			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
//			return err
//		}
//
//		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
//		return err
//	}
//	err = user.DeleteUserEmail()
//	if err != nil {
//		if errors.Is(err, domain.ErrUserIsDeleted) {
//			h.observer.Logger.Trace().Err(err).Msg("user is deleted")
//			return err
//		}
//
//		if errors.Is(err, domain.ErrInvalidEmail) {
//			h.observer.Logger.Error().Err(err).Msg("invalid email")
//			return err
//		}
//	}
//	err = h.userStore.UpdateUser(ctx, user)
//	if err != nil {
//		if errors.Is(err, adapter.ErrStorage) {
//			h.observer.Logger.Error().Err(err).Msg("database error while updating user")
//			return err
//		}
//
//		h.observer.Logger.Error().Err(err).Msg("unexpected error while updating user")
//		return err
//	}
//	return nil
//}
