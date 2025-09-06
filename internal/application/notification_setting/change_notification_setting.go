package notification_setting

import (
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type ChangeNotificationSettingCommand struct {
	ID             string
	FrequencyID    *string
	NotifyTypesIDs *[]string
	MesTypesIDs    *[]string
}

type ChangeNotificationSettingHandler struct {
	notificationSettingStorage db.NotificationSettingRepository
	frequencyStorage           db.NotificationFrequencyRepository
	//notificationTypeStorage    db.NotificationTypeRepository
	//notifySettingNotifyTypeStorage db.NotifySettingNotifyTypeRepository
	//notifySettingMesTypeStorage    db.NotifySettingMesTypeRepository
	//messengerTypeStorage           db.MessengerTypeRepository
	membershipStorage db.FamilyMembershipRepository
	userStorage       db.UserRepository
	familyStorage     db.FamilyRepository
	roleStorage       db.RoleRepository
	zoneStorage       db.ZoneRepository
	observer          *observability.Observability
}

func NewChangeNotificationSettingHandler(
	notificationSettingStorage db.NotificationSettingRepository,
	frequencyStorage db.NotificationFrequencyRepository,
	//notificationTypeStorage db.NotificationTypeRepository,
	//notifySettingNotifyTypeStorage db.NotifySettingNotifyTypeRepository,
	//notifySettingMesTypeStorage db.NotifySettingMesTypeRepository,
	//messengerTypeStorage db.MessengerTypeRepository,
	membershipStorage db.FamilyMembershipRepository,
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	roleStorage db.RoleRepository,
	zoneStorage db.ZoneRepository,
	observer *observability.Observability,
) *ChangeNotificationSettingHandler {
	return &ChangeNotificationSettingHandler{
		notificationSettingStorage: notificationSettingStorage,
		frequencyStorage:           frequencyStorage,
		//notificationTypeStorage:    notificationTypeStorage,
		//notifySettingNotifyTypeStorage: notifySettingNotifyTypeStorage,
		//notifySettingMesTypeStorage:    notifySettingMesTypeStorage,
		//messengerTypeStorage:           messengerTypeStorage,
		membershipStorage: membershipStorage,
		userStorage:       userStorage,
		familyStorage:     familyStorage,
		roleStorage:       roleStorage,
		zoneStorage:       zoneStorage,
		observer:          observer,
	}
}

func (h *ChangeNotificationSettingHandler) Handle(
	ctx context.Context,
	cmd ChangeNotificationSettingCommand,
) error {
	//settingID, err := value_object.NewIDFromString(cmd.ID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid setting id")
	//	return fmt.Errorf("invalid setting id: %w", err)
	//}
	//setting, err := h.notificationSettingStorage.GetNotificationSetting(ctx, settingID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get setting")
	//	return fmt.Errorf("failed to get setting: %w", err)
	//}
	//frequencyID, err := value_object.NewIDFromString(setting.FrequencyID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid frequency id")
	//	return fmt.Errorf("invalid frequency id: %w", err)
	//}
	//frequency, err := h.GetFrequency(ctx, frequencyID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get frequency")
	//	return fmt.Errorf("failed to get frequency: %w", err)
	//}
	//notifyTypes, err := h.notificationTypeStorage.GetNotificationTypesBySetting(ctx, settingID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get notification types by setting")
	//	return fmt.Errorf("failed to get notification types by setting: %w", err)
	//}
	////mesTypes, err := h.messengerTypeStorage.GetMessengerTypesBySetting(ctx, settingID)
	////if err != nil {
	////	h.observer.Logger.Error().Err(err).Msg("failed to get messenger types by setting")
	////	return fmt.Errorf("failed to get messenger types by setting: %w", err)
	////}
	//
	//receiverID, err := value_object.NewIDFromString(setting.ReceiverID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid receiver id")
	//	return fmt.Errorf("invalid receiver id: %w", err)
	//}
	//senderID, err := value_object.NewIDFromString(setting.SenderID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid sender id")
	//	return fmt.Errorf("invalid sender id: %w", err)
	//}
	//users, err := h.userStorage.GetUsersByIDs(ctx, []string{receiverID.String(), senderID.String()})
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get sender and receiver")
	//	return fmt.Errorf("failed to get sender and receiver: %w", err)
	//}
	//
	//zoneID, err := value_object.NewIDFromString(setting.ZoneID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
	//	return fmt.Errorf("invalid zone id: %w", err)
	//}
	//zone, err := h.zoneStorage.GetZone(ctx, zoneID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get zone")
	//	return fmt.Errorf("failed to get zone: %w", err)
	//}
	//
	//familyID, err := value_object.NewIDFromString(setting.FamilyID)
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid family id")
	//	return fmt.Errorf("invalid family id: %w", err)
	//}
	//family, err := h.familyStorage.GetFamily(ctx, familyID)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get family")
	//	return fmt.Errorf("failed to get family: %w", err)
	//}
	//mesContacts, err := users[1].GetUserContactsByMessengerType(mesTypes)
	//if err != nil {
	//	h.observer.Logger.Error().Err(err).Msg("failed to get messenger contacts")
	//	return fmt.Errorf("failed to get user messenger contacts: %w", err)
	//}
	//settingAggregate := &aggregate.NotificationSetting{
	//	ID:                settingID,
	//	Frequency:         frequency,
	//	EventType:         notifyTypes,
	//	MessengerTypes:    mesTypes,
	//	MessengerContacts: mesContacts,
	//	Receiver:          users[1],
	//	Sender:            users[0],
	//	AlertsHistory:     nil,
	//	CreatedAt:         setting.CreatedAt,
	//	UpdatedAt:         setting.UpdatedAt,
	//	Zone:              zone,
	//	Family:            family,
	//}
	//
	//if cmd.FrequencyID != nil {
	//	frequencyID, err = value_object.NewIDFromString(*cmd.FrequencyID)
	//	if err != nil {
	//		h.observer.Logger.Trace().Err(err).Msg("invalid frequency id")
	//		return fmt.Errorf("invalid frequency id: %w", err)
	//	}
	//	frequency, err = h.GetFrequency(ctx, frequencyID)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to get frequency")
	//		return fmt.Errorf("failed to get frequency: %w", err)
	//	}
	//	err = settingAggregate.ChangeFrequency(frequency)
	//	if err != nil {
	//		h.observer.Logger.Trace().Err(err).Msg("failed to change frequency")
	//		return fmt.Errorf("failed to change frequency: %w", err)
	//	}
	//	err = h.notificationSettingStorage.UpdateNotificationSetting(ctx, settingAggregate)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to update notification setting")
	//		return fmt.Errorf("failed to update notification setting: %w", err)
	//	}
	//}
	//
	//if cmd.NotifyTypesIDs != nil {
	//	var notifyTypesIDs []value_object.ID
	//	for _, id := range *cmd.NotifyTypesIDs {
	//		notifyTypesID, err := value_object.NewIDFromString(id)
	//		if err != nil {
	//			h.observer.Logger.Trace().Err(err).Msg("invalid notification type id")
	//			return fmt.Errorf("invalid notification type id: %w", err)
	//		}
	//		notifyTypesIDs = append(notifyTypesIDs, notifyTypesID)
	//	}
	//	notifyTypes, err = h.notificationTypeStorage.GetNotificationTypesByIDs(ctx, *cmd.NotifyTypesIDs)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to get notification types")
	//		return fmt.Errorf("failed to get notification types: %w", err)
	//	}
	//	err = settingAggregate.ChangeNotificationTypes(notifyTypes)
	//	if err != nil {
	//		h.observer.Logger.Trace().Err(err).Msg("failed to change notification types")
	//		return fmt.Errorf("failed to change notification types: %w", err)
	//	}
	//
	//	// изменить на update
	//	//err = h.notifySettingNotifyTypeStorage.DeleteNotifySettingNotifyTypesByNotifySetting(ctx, settingAggregate.externalID)
	//	//if err != nil {
	//	//	h.observer.Logger.Error().Err(err).Msg("failed to delete notify setting and notify types relationships")
	//	//	return fmt.Errorf("failed to delete notify setting and notify types relationships: %w", err)
	//	//}
	//
	//	var notifySettingNotTypes []*entity.NotifySettingNotifyType
	//	for _, notifyType := range notifyTypes {
	//		notSettingNotType, err := entity.NewNotifySettingNotifyType(settingAggregate.ID, notifyType.ID)
	//		if err != nil {
	//			h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting and notify type relation")
	//			return fmt.Errorf("failed to create notify setting and notify type relation: %w", err)
	//		}
	//		notifySettingNotTypes = append(notifySettingNotTypes, notSettingNotType)
	//	}
	//
	//	err = h.notifySettingNotifyTypeStorage.AddNotifySettingNotTypes(ctx, notifySettingNotTypes)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to add notify setting and notify type relation to db")
	//		return fmt.Errorf("failed to add notify setting and notify type relation to db: %w", err)
	//	}
	//}
	//
	//if cmd.MesTypesIDs != nil {
	//	var mesTypesIDs []value_object.ID
	//	for _, id := range *cmd.MesTypesIDs {
	//		mesTypesID, err := value_object.NewIDFromString(id)
	//		if err != nil {
	//			h.observer.Logger.Trace().Err(err).Msg("invalid messenger type id")
	//			return fmt.Errorf("invalid messenger type id: %w", err)
	//		}
	//		mesTypesIDs = append(mesTypesIDs, mesTypesID)
	//	}
	//	mesTypes, err = h.messengerTypeStorage.GetMessengerTypesByIDs(ctx, *cmd.MesTypesIDs)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to get messenger types")
	//		return fmt.Errorf("failed to get messenger types: %w", err)
	//	}
	//	err = settingAggregate.ChangeMessengerTypes(mesTypes)
	//	if err != nil {
	//		h.observer.Logger.Trace().Err(err).Msg("failed to change messenger types")
	//		return fmt.Errorf("failed to change messenger types: %w", err)
	//	}
	//
	//	/// Переделать
	//	err = h.notifySettingMesTypeStorage.DeleteNotifySettingMesTypesByNotifySetting(ctx, settingAggregate.ID)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to delete notify setting and mes types relationships")
	//		return fmt.Errorf("failed to delete notify setting and mes types relationships: %w", err)
	//	}
	//
	//	var notifySettingMesTypes []*entity.NotifySettingMesType
	//	for _, mesType := range mesTypes {
	//		notSettingMesType, err := entity.NewNotifySettingMesType(settingAggregate.ID, mesType.ID)
	//		if err != nil {
	//			h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting and messenger type relation")
	//			return fmt.Errorf("failed to create notify setting and messenger type relation: %w", err)
	//		}
	//		notifySettingMesTypes = append(notifySettingMesTypes, notSettingMesType)
	//	}
	//	err = h.notifySettingMesTypeStorage.AddNotifySettingMesTypes(ctx, notifySettingMesTypes)
	//	if err != nil {
	//		h.observer.Logger.Error().Err(err).Msg("failed to add notify setting and messenger type relation to db")
	//		return fmt.Errorf("failed to add notify setting and messenger type relation to db: %w", err)
	//	}
	//}
	return nil
}

//func (h *ChangeNotificationSettingHandler) GetFrequency(ctx context.Context,
//	id value_object.ID) (*entity.NotificationFrequency, error) {
//	freq, err := h.frequencyStorage.GetFrequency(ctx, id)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get frequency: %w", err)
//	}
//	return freq, nil
//}
