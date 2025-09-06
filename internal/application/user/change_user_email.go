package user

//type ChangeUserEmailCommand struct {
//	ID    uuid.UUID
//	Email string
//}
//
//type ChangeUserEmailHandler struct {
//	storage  db.UserRepository
//	observer *observability.Observability
//}
//
//func NewChangeUserEmailHandler(
//	storage db.UserRepository,
//	observer *observability.Observability,
//) *ChangeUserEmailHandler {
//	return &ChangeUserEmailHandler{
//		storage:  storage,
//		observer: observer,
//	}
//}
//
//func (h *ChangeUserEmailHandler) Handle(ctx context.Context, cmd ChangeUserEmailCommand) error {
//	email, err := value_object.NewEmail(cmd.Email)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to add email")
//		return fmt.Errorf("failed to add email: %w", err)
//	}
//	exists, err := h.storage.EmailExists(ctx, *email)
//	if err != nil {
//		if errors.Is(err, adapter.ErrStorage) {
//			h.observer.Logger.Error().Err(err).Msg("database error while getting email exists")
//			return err
//		}
//
//		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting email exists")
//		return err
//	}
//	if exists {
//		h.observer.Logger.Trace().Msg("this email already exists in system")
//		return domain.ErrEmailExists
//	}
//
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
//	err = user.ChangeEmail(*email)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to change email")
//		return fmt.Errorf("failed to change email: %w", err)
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
//	return nil
//}
