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
	"github.com/google/uuid"
)

type GetNotificationSettingsQuery struct {
	UserID string
	ZoneID uuid.UUID
}

type GetNotificationSettingsHandler struct {
	notificationSettingStorage db.NotificationSettingRepository
	frequencyStorage           db.NotificationFrequencyRepository
	userStorage                db.UserRepository
	familyStorage              db.FamilyRepository
	zoneStorage                db.ZoneRepository
	//notSettingNotTypeStorage   db.NotifySettingNotifyTypeRepository
	//notSettingMesTypeStorage   db.NotifySettingMesTypeRepository
	observer *observability.Observability
}

func NewGetNotificationSettingsHandler(
	notificationSettingStorage db.NotificationSettingRepository,
	frequencyStorage db.NotificationFrequencyRepository,
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	zoneStorage db.ZoneRepository,
	//notSettingNotTypeStorage db.NotifySettingNotifyTypeRepository,
	//notSettingMesTypeStorage db.NotifySettingMesTypeRepository,
	observer *observability.Observability,
) *GetNotificationSettingsHandler {
	return &GetNotificationSettingsHandler{
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

func (h *GetNotificationSettingsHandler) Handle(ctx context.Context,
	query GetNotificationSettingsQuery) ([]*app_model.ApplicationNotificationSettings, error) {
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

	zoneID, err := value_object.NewIDFromString(query.ZoneID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
		return nil, err
	}
	zone, err := h.zoneStorage.GetZone(ctx, zoneID)
	if err != nil {
		if errors.Is(err, adapter.ErrZoneNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("zone not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting zone")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone")
		return nil, err
	}
	fmt.Printf("zone: %v\n", zone)

	dbSettings, err := h.notificationSettingStorage.GetNotificationSettingsByZoneForUser(ctx, receiver.ID, zone.ID)
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

	var freqIDs, senderIDs []string
	for _, ns := range dbSettings {
		if ns.FrequencyID != nil {
			freqIDs = append(freqIDs, *ns.FrequencyID)
		}
		senderIDs = append(senderIDs, ns.SenderID)
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
			h.observer.Logger.Trace().Err(err).Msg("sender not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notify senders")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify senders")
		return nil, err
	}

	appSettings, err := app_model.NewApplicationNotificationSettingsList(dbSettings, frequencies, senders)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("failed to map application settings")
		return nil, err
	}
	return appSettings, nil

	//freqMap := make(map[string]*entity.NotificationFrequency, len(frequencies))
	//for _, f := range frequencies {
	//	freqMap[f.ID.String()] = f
	//}
	//senderMap := make(map[string]*entity.User, len(senders))
	//for _, u := range senders {
	//	senderMap[u.ID.String()] = u
	//}
	//
	//var result []*app_model.ApplicationNotificationSettings
	//for _, ns := range dbSettings {
	//	voID, err := value_object.NewIDFromString(ns.ID)
	//	if err != nil {
	//		return nil, fmt.Errorf("invalid setting id %q: %w", ns.ID, err)
	//	}
	//	var appFreq *app_model.ApplicationNotificationFrequency
	//	if ns.FrequencyID != nil {
	//		entFreq, ok := freqMap[*ns.FrequencyID]
	//		if !ok {
	//			return nil, fmt.Errorf("frequency %q not found", *ns.FrequencyID)
	//		}
	//		appFreq = app_model.NewApplicationNotificationFrequency(entFreq)
	//	}
	//	entUser, ok := senderMap[ns.SenderID]
	//	if !ok {
	//		return nil, fmt.Errorf("sender %q not found", ns.SenderID)
	//	}
	//	appSender := app_model.NewApplicationUser(entUser)
	//
	//	var minBattery *int
	//	setting, err := app_model.NewApplicationNotificationSetting(
	//		voID.ToRaw(),
	//		appFreq,
	//		appSender,
	//		minBattery,
	//	)
	//	result = append(result, setting)
	//}
	//
	//return result, nil
}

//notifySettings, err := h.notificationSettingStorage.GetNotificationSettingsByZoneForUser(ctx, receiver.ID, zone.ID)
//if err != nil {
//	if errors.Is(err, adapter.ErrStorage) {
//		h.observer.Logger.Error().Err(err).Msg("database error while getting notify settings")
//		return nil, err
//	}
//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify settings")
//	return nil, err
//}
//if len(notifySettings) == 0 {
//	h.observer.Logger.Trace().Msg("no notify settings found")
//	return nil, nil
//}
//fmt.Printf("notify settings: %v\n", notifySettings)

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
