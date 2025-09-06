package notification_setting

import (
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"fmt"
)

type DeleteNotificationSettingCommand struct {
	ID string
}

type DeleteNotificationSettingHandler struct {
	notificationSettingStorage db.NotificationSettingRepository
	//notifySettingNotifyTypeStorage db.NotifySettingNotifyTypeRepository
	//notifySettingMesTypeStorage    db.NotifySettingMesTypeRepository
	observer *observability.Observability
}

func NewDeleteNotificationSettingHandler(
	notificationSettingStorage db.NotificationSettingRepository,
	//notifySettingNotifyTypeStorage db.NotifySettingNotifyTypeRepository,
	//notifySettingMesTypeStorage db.NotifySettingMesTypeRepository,
	observer *observability.Observability,
) *DeleteNotificationSettingHandler {
	return &DeleteNotificationSettingHandler{
		notificationSettingStorage: notificationSettingStorage,
		//notifySettingNotifyTypeStorage: notifySettingNotifyTypeStorage,
		//notifySettingMesTypeStorage:    notifySettingMesTypeStorage,
		observer: observer,
	}
}

func (h *DeleteNotificationSettingHandler) Handle(ctx context.Context, cmd DeleteNotificationSettingCommand) error {
	notifySettingID, err := value_object.NewIDFromString(cmd.ID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid notification setting id")
		return fmt.Errorf("invalid notification setting id: %w", err)
	}
	_, err = h.notificationSettingStorage.GetNotificationSetting(ctx, notifySettingID)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to get notification setting")
		return fmt.Errorf("failed to get notification setting: %w", err)
	}

	/// изменить на update
	//err = h.notifySettingNotifyTypeStorage.DeleteNotifySettingNotifyTypesByNotifySetting(ctx, notifySettingID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to delete notify setting and notify types relationships")
	//	return fmt.Errorf("failed to delete notify setting and notify types relationships: %w", err)
	//}
	// TODO: Переделать
	//err = h.notifySettingMesTypeStorage.DeleteNotifySettingMesTypesByNotifySetting(ctx, notifySettingID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to delete notify setting and mes types relationships")
	//	return fmt.Errorf("failed to delete notify setting and mes types relationships: %w", err)
	//}
	//// Заменить на update
	//err = h.notificationSettingStorage.DeleteNotificationSetting(ctx, notifySettingID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to delete notification setting")
	//	return fmt.Errorf("failed to delete notification setting: %w", err)
	//}
	return nil
}
