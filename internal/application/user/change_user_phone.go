package user

//type ChangeUserPhoneCommand struct {
//	ID    uuid.UUID
//	Phone string
//}
//
//type ChangeUserPhoneHandler struct {
//	storage  db.UserRepository
//	observer *observability.Observability
//}
//
//func NewChangeUserPhoneHandler(
//	storage db.UserRepository,
//	observer *observability.Observability,
//) *ChangeUserPhoneHandler {
//	return &ChangeUserPhoneHandler{
//		storage:  storage,
//		observer: observer,
//	}
//}
//
//func (h *ChangeUserPhoneHandler) Handle(ctx context.Context, cmd ChangeUserPhoneCommand) error {
//	phone, err := value_object.NewPhone(cmd.Phone)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to add phone")
//		return fmt.Errorf("failed to add phone: %w", err)
//	}
//	exists, err := h.storage.PhoneExists(ctx, phone)
//	if err != nil {
//		if errors.Is(err, adapter.ErrStorage) {
//			h.observer.Logger.Error().Err(err).Msg("database error while getting phone exists")
//			return err
//		}
//
//		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting phone exists")
//		return err
//	}
//	if exists {
//		h.observer.Logger.Trace().Msg("this phone already exists in system")
//		return domain.ErrPhoneExists
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
//	err = user.ChangePhone(phone)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to change phone")
//		return fmt.Errorf("failed to change phone: %w", err)
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
