package notification_setting

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
)

type GetBatterySettingsQuery struct {
	UserID string
}

type GetBatterySettingsHandler struct {
	notificationSettingStorage db.NotificationSettingRepository
	frequencyStorage           db.NotificationFrequencyRepository
	userStorage                db.UserRepository
	familyStorage              db.FamilyRepository
	zoneStorage                db.ZoneRepository
	//notSettingNotTypeStorage   db.NotifySettingNotifyTypeRepository
	//notSettingMesTypeStorage   db.NotifySettingMesTypeRepository
	observer *observability.Observability
}

func NewGetBatterySettingsHandler(
	notificationSettingStorage db.NotificationSettingRepository,
	frequencyStorage db.NotificationFrequencyRepository,
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	zoneStorage db.ZoneRepository,
	//notSettingNotTypeStorage db.NotifySettingNotifyTypeRepository,
	//notSettingMesTypeStorage db.NotifySettingMesTypeRepository,
	observer *observability.Observability,
) *GetBatterySettingsHandler {
	return &GetBatterySettingsHandler{
		notificationSettingStorage: notificationSettingStorage,
		frequencyStorage:           frequencyStorage,
		userStorage:                userStorage,
		familyStorage:              familyStorage,
		zoneStorage:                zoneStorage,
		//notSettingNotTypeStorage:   notSettingNotTypeStorage,
		//notSettingMesTypeStorage:   notSettingMesTypeStorage,
		observer: observer,
	}
}

func (h *GetBatterySettingsHandler) Handle(ctx context.Context,
	query GetBatterySettingsQuery) ([]*app_model.ApplicationNotificationSettings, error) {
	extID, err := value_object.NewIDFromString(query.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return nil, err
	}
	receiver, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("notify receiver not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notify receiver")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify receiver")
		return nil, err
	}
	fmt.Printf("receiver: %v\n", receiver)

	dbSettings, err := h.notificationSettingStorage.GetBatteryNotificationSettingsByReceiver(ctx, receiver.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notify settings")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify settings")
		return nil, err
	}
	if len(dbSettings) == 0 {
		h.observer.Logger.Trace().Msg("no notify settings found")
		return nil, nil
	}
	fmt.Printf("notify settings: %v\n", dbSettings)

	freqIDs := make([]string, 0, len(dbSettings))
	senderIDs := make([]string, 0, len(dbSettings))
	for _, s := range dbSettings {
		freqIDs = append(freqIDs, *s.FrequencyID)
		senderIDs = append(senderIDs, s.SenderID)
	}

	frequencies, err := h.frequencyStorage.GetNotificationFrequenciesByIDs(ctx, freqIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notify frequency")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify frequency")
		return nil, err
	}

	senders, err := h.userStorage.GetUsersByIDs(ctx, senderIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return nil, err
	}

	return app_model.NewApplicationNotificationSettingsList(
		dbSettings,
		frequencies,
		senders,
	)

	//childMap := make(map[string]*entity.User, len(children))
	//for _, c := range children {
	//	childMap[c.ID.String()] = c
	//}
	//
	//var result []*app_model.ApplicationNotificationSettings
	//for _, s := range dbSettings {
	//	settingIDVo, err := value_object.NewIDFromString(s.NotifySettingID)
	//	if err != nil {
	//		h.observer.Logger.Trace().Err(err).Msg("invalid notify_setting_id")
	//		return nil, fmt.Errorf("invalid notification_id %q: %w", s.NotifySettingID, err)
	//	}
	//	settingUUID := settingIDVo.ToRaw()
	//
	//	var appFreq *app_model.ApplicationNotificationFrequency
	//	if s.FrequencyID != "" {
	//		freqIDVo, err := value_object.NewIDFromString(s.FrequencyID)
	//		if err != nil {
	//			h.observer.Logger.Trace().Err(err).Msg("invalid frequency_id")
	//			return nil, fmt.Errorf("invalid frequency_id %q: %w", s.FrequencyID, err)
	//		}
	//		appFreq = &app_model.ApplicationNotificationFrequency{
	//			ID:        freqIDVo.ToRaw(),
	//			Frequency: s.Frequency,
	//		}
	//	}
	//
	//	childEntity, ok := childMap[s.SenderID]
	//	if !ok {
	//		h.observer.Logger.Trace().Msgf("child with id %s not found", s.SenderID)
	//		return nil, fmt.Errorf("child with id %q not found", s.SenderID)
	//	}
	//	appChild := app_model.NewApplicationUser(childEntity)
	//	minBat := s.Threshold
	//
	//	result = append(result, &app_model.ApplicationNotificationSettings{
	//		ID:         settingUUID,
	//		Frequency:  appFreq,
	//		Sender:     appChild,
	//		MinBattery: &minBat,
	//	})
	//}
	//
	//return result, nil
}

//children, err := h.userStorage.GetUsersByIDs(ctx, childIDs)
//
//childMap := make(map[string]*entity.User, len(children))
//for _, c := range children {
//	childMap[c.ID.String()] = c
//}
//
//details := make([]*app_model.NotificationSettingDetails, 0, len(dbSettings))
//for _, s := range dbSettings {
//	minBat := s.Threshold
//	settingUUID, err := value_object.NewIDFromString(s.NotifySettingID)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("invalid notify setting id")
//		return nil, fmt.Errorf("invalid notification_id %q: %w", s.NotifySettingID, err)
//	}
//	freqUUID, err := value_object.NewIDFromString(s.FrequencyID)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("invalid frequency id")
//		return nil, fmt.Errorf("invalid frequency_id %q: %w", s.FrequencyID, err)
//	}
//	appFreq := &app_model.ApplicationNotificationFrequency{
//		ID:        freqUUID.ToRaw(),
//		Frequency: s.Frequency,
//	}
//
//	childEntity, ok := childMap[s.SenderID]
//	if !ok {
//		h.observer.Logger.Trace().Err(err).Msg("child not found")
//		return nil, fmt.Errorf("child with id %q not found", s.SenderID)
//	}
//	appChild := app_model.NewApplicationUser(childEntity)
//
//	details = append(details, &app_model.NotificationSettingDetails{
//		ID:         settingUUID.ToRaw(),
//		Frequency:  appFreq,
//		Sender:     appChild,
//		MinBattery: &minBat,
//	})
//}
//
//return &app_model.ApplicationNotificationSettings{
//	EventType: "battery",
//	Details:   details,
//}, nil

//
//var frequencyIDs []string
//for _, ns := range notifySettings {
//	if ns.FrequencyID != nil {
//		frequencyIDs = append(frequencyIDs, *ns.FrequencyID)
//	}
//}
//fmt.Printf("frequency ids: %v", frequencyIDs)
//
//var frequencies []*entity.NotificationFrequency
//if len(frequencyIDs) > 0 {
//	fmt.Printf("frequencies > 0")
//	frequencies, err = h.frequencyStorage.GetNotificationFrequenciesByIDs(ctx, frequencyIDs)
//	if err != nil {
//		if errors.Is(err, adapter.ErrStorage) {
//			h.observer.Logger.Error().Err(err).Msg("database error while getting notify frequency")
//			return nil, err
//		}
//		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify frequency")
//		return nil, err
//	}
//} else {
//	frequencies = []*entity.NotificationFrequency{}
//	fmt.Printf("frequencies < 0")
//}
//
//sendersIDS := make([]string, len(notifySettings))
//for key, setting := range notifySettings {
//	sendersIDS[key] = setting.SenderID
//}
//fmt.Printf("senders: %v\n", sendersIDS)
//senders, err := h.userStorage.GetUsersByIDs(ctx, sendersIDS)
//if err != nil {
//	if errors.Is(err, adapter.ErrUserNotFound) {
//		h.observer.Logger.Trace().Err(err).Msg("sender not found")
//		return nil, err
//	}
//	if errors.Is(err, adapter.ErrStorage) {
//		h.observer.Logger.Error().Err(err).Msg("database error while getting notify senders")
//		return nil, err
//	}
//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify senders")
//	return nil, err
//}
//
//resSettings, err := app_model.NewApplicationNotificationSettings(notifySettings, frequencies, senders)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to create application notification settings")
//	return nil, fmt.Errorf("failed to create application notification settings: %w", err)
//}
//
//return resSettings, nil

//userID, err := value_object.NewIDFromString(query.UserID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid user id")
//	return nil, fmt.Errorf("invalid user id: %w", err)
//}
//zoneID, err := value_object.NewIDFromString(query.ZoneID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
//	return nil, fmt.Errorf("invalid zone id: %w", err)
//}
//zone, err := h.zoneStorage.GetZone(ctx, zoneID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get zone")
//	return nil, fmt.Errorf("failed to get zone: %w", err)
//}
//familyID, err := value_object.NewIDFromString(query.FamilyID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid family id")
//	return nil, fmt.Errorf("invalid family id: %w", err)
//}
//family, err := h.familyStorage.GetFamily(ctx, familyID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get family")
//	return nil, fmt.Errorf("failed to get family: %w", err)
//}
//notificationSettings, err := h.notificationSettingStorage.GetNotificationSettingsByZoneReceiverAndFamily(ctx,
//	userID, familyID, zoneID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get notification settings")
//	return nil, fmt.Errorf("failed to get notification settings: %w", err)
//}
//notificationSettingsIDs := make([]string, len(notificationSettings))
//notificationFrequenciesIDs := make([]string, len(notificationSettings))
//notificationSendersIDs := make([]string, len(notificationSettings))
//for i, setting := range notificationSettings {
//	notificationSettingsIDs[i] = setting.ID
//	notificationFrequenciesIDs[i] = setting.FrequencyID
//	notificationSendersIDs[i] = setting.SenderID
//}
//notificationFrequencies, err := h.frequencyStorage.GetNotificationFrequenciesByIDs(ctx, notificationFrequenciesIDs)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get notification frequencies")
//	return nil, fmt.Errorf("failed to get notification frequencies: %w", err)
//}
//notificationSenders, err := h.userStorage.GetUsersByIDs(ctx, notificationSendersIDs)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get notification senders")
//	return nil, fmt.Errorf("failed to get notification senders: %w", err)
//}
//
//notifySettingsNotifyTypes, err := h.notSettingNotTypeStorage.GetNotifySettingsNotifyTypesBySettings(ctx,
//	notificationSettingsIDs)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get notify settings and notify types relations")
//	return nil, fmt.Errorf("failed to get notify settings and notify types relations: %w", err)
//}
//notifySettingsMesTypes, err := h.notSettingMesTypeStorage.GetNotifySettingsMesTypesBySettings(ctx,
//	notificationSettingsIDs)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get notify settings and mes types relations")
//	return nil, fmt.Errorf("failed to get notify settings and mes types relations: %w", err)
//}
//
