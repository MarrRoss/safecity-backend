package notification_setting

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type GetNotificationsSendersByReceiverCommand struct {
	ID string
}

type GetNotificationsSendersByReceiverHandler struct {
	notificationSettingStorage db.NotificationSettingRepository
	userStorage                db.UserRepository
}

func NewGetNotificationsSendersByReceiverHandler(
	notificationSettingStorage db.NotificationSettingRepository,
	userStorage db.UserRepository) *GetNotificationsSendersByReceiverHandler {
	return &GetNotificationsSendersByReceiverHandler{
		notificationSettingStorage: notificationSettingStorage,
		userStorage:                userStorage,
	}
}

func (h *GetNotificationsSendersByReceiverHandler) Handle(ctx context.Context,
	cmd GetNotificationsSendersByReceiverCommand) (*app_model.ApplicationNotificationsSendersByReceiver, error) {
	//receiverID, err := value_object.NewIDFromString(cmd.externalID)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid receiver id: %w", err)
	//}
	//user, err := h.userStorage.GetUser(ctx, receiverID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get user: %w", err)
	//}
	//if user.EndedAt != nil {
	//	return nil, fmt.Errorf("receiver user has been deleted")
	//}
	//notifSettingsIDs, err := h.notificationSettingStorage.GetNotificationsSettingsIDsByReceiver(ctx, user.externalID)
	//if err != nil {
	//	return nil, err
	//}
	//var notifSettings []*response.NotificationSettingDB
	//for _, notifSettingID := range notifSettingsIDs {
	//	notifSetting, err := h.notificationSettingStorage.GetNotificationSetting(ctx, notifSettingID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	notifSettings = append(notifSettings, notifSetting)
	//}
	//var senderIDs []value_object.externalID
	//for _, notifSetting := range notifSettings {
	//	if notifSetting.SenderID == "" {
	//		return nil, fmt.Errorf("sender id is nil")
	//	}
	//	id, err := value_object.NewIDFromString(notifSetting.SenderID)
	//	if err != nil {
	//		return nil, fmt.Errorf("invalid sender id: %w", err)
	//	}
	//	senderIDs = append(senderIDs, id)
	//}
	//var senders map[string]string
	//for _, senderID := range senderIDs {
	//	sender, err := h.userStorage.GetUser(ctx, senderID)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to get sender: %w", err)
	//	}
	//	if sender.EndedAt != nil {
	//		return nil, fmt.Errorf("sender user has been deleted")
	//	}
	//	senders[sender.Phone.String()] = sender.Name.FirstName.String() + " " +
	//		sender.Name.LastName.String() + " " + sender.Name.Patronymic.String()
	//}
	//
	//return app_model.NewApplicationNotificationsSendersByReceiver(senders), nil

	return nil, nil
}
