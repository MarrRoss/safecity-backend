package user

//type ChangeUserPasswordCommand struct {
//	ID       uuid.UUID
//	Password string
//}
//
//type ChangeUserPasswordHandler struct {
//	storage  db.UserRepository
//	observer *observability.Observability
//}
//
//func NewChangeUserPasswordHandler(
//	storage db.UserRepository,
//	observer *observability.Observability,
//) *ChangeUserPasswordHandler {
//	return &ChangeUserPasswordHandler{
//		storage:  storage,
//		observer: observer,
//	}
//}
//
//func (h *ChangeUserPasswordHandler) Handle(ctx context.Context, cmd ChangeUserPasswordCommand) error {
//	userID, err := value_object.NewIDFromString(cmd.ID.String())
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("invalid user id")
//		return err
//	}
//	user, err := h.storage.GetUser(ctx, userID)
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
//
//	password, err := value_object.NewPassword(cmd.Password)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to add password")
//		return fmt.Errorf("failed to add password: %w", err)
//	}
//	err = user.ChangePassword(password)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to change password")
//		return fmt.Errorf("failed to change password: %w", err)
//	}
//	err = h.storage.UpdateUser(ctx, user)
//	if err != nil {
//		if errors.Is(err, adapter.ErrStorage) {
//			h.observer.Logger.Error().Err(err).Msg("database error while updating user")
//			return err
//		}
//
//		h.observer.Logger.Error().Err(err).Msg("unexpected error while updating user")
//		return err
//	}
//
//	return nil
//}
