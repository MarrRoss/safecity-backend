package app_model

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"github.com/google/uuid"
)

type ApplicationNotificationSettings struct {
	ID         uuid.UUID
	Frequency  *ApplicationNotificationFrequency
	Sender     *ApplicationUser
	MinBattery *int
	//NotificationTypes []*ApplicationNotificationType
	//MessengerTypes    []*ApplicationMessengerType
}

func NewApplicationNotificationSetting(
	id uuid.UUID,
	frequency *ApplicationNotificationFrequency,
	sender *ApplicationUser,
	minBattery *int,
) (*ApplicationNotificationSettings, error) {
	return &ApplicationNotificationSettings{
		ID:         id,
		Frequency:  frequency,
		Sender:     sender,
		MinBattery: minBattery,
	}, nil
}

func NewApplicationNotificationSettingsList(
	dbSettings []*response.NotificationSettingMinBatteryDB,
	frequencies []*entity.NotificationFrequency,
	senders []*entity.User,
) ([]*ApplicationNotificationSettings, error) {
	freqMap := make(map[string]*entity.NotificationFrequency, len(frequencies))
	for _, f := range frequencies {
		freqMap[f.ID.String()] = f
	}
	senderMap := make(map[string]*entity.User, len(senders))
	for _, u := range senders {
		senderMap[u.ID.String()] = u
	}

	result := make([]*ApplicationNotificationSettings, 0, len(dbSettings))
	for _, ns := range dbSettings {
		voID, err := value_object.NewIDFromString(ns.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid setting id %q: %w", ns.ID, err)
		}

		var appFreq *ApplicationNotificationFrequency
		if ns.FrequencyID != nil {
			entFreq, ok := freqMap[*ns.FrequencyID]
			if !ok {
				return nil, fmt.Errorf("frequency %q not found", *ns.FrequencyID)
			}
			appFreq = NewApplicationNotificationFrequency(entFreq)
		}

		entUser, ok := senderMap[ns.SenderID]
		if !ok {
			return nil, fmt.Errorf("sender %q not found", ns.SenderID)
		}
		appSender := NewApplicationUser(entUser)

		var minBat *int
		if ns.MinBattery != nil {
			minBat = ns.MinBattery
		}

		setting, err := NewApplicationNotificationSetting(
			voID.ToRaw(),
			appFreq,
			appSender,
			minBat,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, setting)
	}
	return result, nil
}

//func NewApplicationNotificationSettings(
//	notifySettings []*response.NotificationSettingDB,
//	frequencies []*entity.NotificationFrequency,
//	senders []*entity.User,
//) (*ApplicationNotificationSettings, error) {
//	fmt.Printf("App func")
//	if len(senders) == 0 {
//		return nil, fmt.Errorf("senders list is empty: %w", application.ErrNotifySendersNotFound)
//	}
//
//	freqMap := make(map[string]*entity.NotificationFrequency, len(frequencies))
//	for _, f := range frequencies {
//		freqMap[f.ID.String()] = f
//	}
//	fmt.Printf("freqMap")
//
//	senderMap := make(map[string]*entity.User)
//	for _, s := range senders {
//		senderMap[s.ID.String()] = s
//	}
//	fmt.Printf("senderMap")
//
//	eventType := notifySettings[0].EventType
//	fmt.Printf("eventType")
//	details := make([]*NotificationSettingDetails, 0, len(notifySettings))
//	for _, ns := range notifySettings {
//		settingID, err := value_object.NewIDFromString(ns.ID)
//		if err != nil {
//			return nil, fmt.Errorf("notification setting %s: %w", ns.ID, domain.ErrInvalidID)
//		}
//		fmt.Printf("settingID %v", settingID)
//
//		sender, ok := senderMap[ns.SenderID]
//		if !ok {
//			return nil, fmt.Errorf("sender with id %s not found", ns.SenderID)
//		}
//		fmt.Printf("sender %v", sender)
//		appSender := NewApplicationUser(sender)
//
//		var appFreq *ApplicationNotificationFrequency
//		if ns.FrequencyID != nil {
//			f, ok := freqMap[*ns.FrequencyID]
//			if !ok {
//				return nil, fmt.Errorf("frequency %q not found", *ns.FrequencyID)
//			}
//			appFreq = NewApplicationNotificationFrequency(f)
//		}
//
//		detail, err := NewApplicationNotificationSetting(settingID.ToRaw(), appFreq, appSender, nil)
//		if err != nil {
//			return nil, fmt.Errorf("failed to create notification setting detail: %w", application.ErrInvalidNotifySettingDetail)
//		}
//		details = append(details, detail)
//	}
//
//	result := &ApplicationNotificationSettings{
//		EventType: eventType,
//		Details:   details,
//	}
//	return result, nil
//}

//func NewApplicationNotificationSettings(
//	notifySettings []*response.NotificationSettingDB,
//	frequencies []*entity.NotificationFrequency,
//	senders []*entity.User,
//) (*ApplicationNotificationSettings, error) {
//	if len(notifySettings) != len(frequencies) || len(frequencies) != len(senders) {
//		return nil, fmt.Errorf("lengths of notifySettings, frequencies and senders do not match: %v", application.ErrDifferentArrayLength)
//	}
//	if len(notifySettings) == 0 {
//		return nil, fmt.Errorf("notify settings list is empty: %v", application.ErrNotifySettingsNotFound)
//	}
//	if len(frequencies) == 0 {
//		return nil, fmt.Errorf("frequencies list is empty: %v", application.ErrNotifyFrequenciesNotFound)
//	}
//	if len(senders) == 0 {
//		return nil, fmt.Errorf("senders list is empty: %v", application.ErrNotifySendersNotFound)
//	}
//	eventType := notifySettings[0].EventType
//	details := make([]*NotificationSettingDetails, len(notifySettings))
//	for i, ns := range notifySettings {
//		settingID, err := value_object.NewIDFromString(ns.ID)
//		if err != nil {
//			return nil, fmt.Errorf("notification setting %s: %w", ns.ID, domain.ErrInvalidID)
//		}
//		appFreq := NewApplicationNotificationFrequency(frequencies[i])
//		appSender := NewApplicationUser(senders[i])
//		detail, err := NewApplicationNotificationSetting(settingID.ToRaw(), appFreq, appSender)
//		if err != nil {
//			return nil, fmt.Errorf("failed to create notification setting detail: %w", application.ErrInvalidNotifySettingDetail)
//		}
//		details[i] = detail
//	}
//
//	result := &ApplicationNotificationSettings{
//		EventType: eventType,
//		Details:   details,
//	}
//	return result, nil
//}

//func NewApplicationNotificationSetting(
//	notificationSetting *response.NotificationSettingDB,
//	frequency *entity.NotificationFrequency,
//	sender *entity.User,
//	//notificationTypes []*ApplicationNotificationType,
//	//messengerTypes []*ApplicationMessengerType,
//	family *entity.Family,
//	eventType string,
//	zone *entity.Zone,
//) *NotificationSettingDetails {
//	return &NotificationSettingDetails{
//		ID:        notificationSetting.ID,
//		Frequency: NewApplicationNotificationFrequency(frequency),
//		Sender:    NewApplicationUser(sender),
//		//NotificationTypes: notificationTypes,
//		//MessengerTypes:    messengerTypes,
//		EventType: eventType,
//		Zone:      NewApplicationZone(zone),
//	}
//}

//func NewApplicationNotificationSettings(
//	notificationSettings []*response.NotificationSettingDB,
//	notificationFrequencies []*entity.NotificationFrequency,
//	notificationSenders []*entity.User,
//	//notifySettingsNotifyTypes []*response.NotificationSettingNotificationTypeDB,
//	//notifySettingsMesTypes []*response.NotificationSettingMessengerTypeDB,
//	eventType string,
//	zone *entity.Zone,
//) ([]*NotificationSettingDetails, error) {
//	frequencyMap := make(map[string]*entity.NotificationFrequency)
//	for _, frequency := range notificationFrequencies {
//		frequencyMap[frequency.ID.String()] = frequency
//	}
//	senderMap := make(map[string]*entity.User)
//	for _, sender := range notificationSenders {
//		senderMap[sender.ID.String()] = sender
//	}
//notificationTypesMap := make(map[string][]*ApplicationNotificationType)
//for _, notifyType := range notifySettingsNotifyTypes {
//	appNotifyType := &ApplicationNotificationType{
//		ID:   notifyType.NotifyTypeID,
//		Type: notifyType.NotifyType,
//	}
//	notificationTypesMap[notifyType.NotifySettingID] = append(notificationTypesMap[notifyType.NotifySettingID], appNotifyType)
//}
//messengerTypesMap := make(map[string][]*ApplicationMessengerType)
//for _, mesType := range notifySettingsMesTypes {
//	appMesType := &ApplicationMessengerType{
//		ID:   mesType.MessengerTypeID,
//		Type: mesType.MessengerType,
//	}
//	messengerTypesMap[mesType.NotifySettingID] = append(messengerTypesMap[mesType.NotifySettingID], appMesType)
//}

//appNotificationSettings := make([]*NotificationSettingDetails, 0, len(notificationSettings))
//for _, setting := range notificationSettings {
//	frequency, ok := frequencyMap[setting.FrequencyID]
//	if !ok {
//		return nil, fmt.Errorf("frequency not found for FrequencyID: %s", setting.FrequencyID)
//	}
//
//	sender, ok := senderMap[setting.SenderID]
//	if !ok {
//		return nil, fmt.Errorf("sender not found for SenderID: %s", setting.SenderID)
//	}
//
//	//notificationTypes := notificationTypesMap[setting.ID]
//	//messengerTypes := messengerTypesMap[setting.ID]
//
//	appNotificationSetting := NewApplicationNotificationSetting(setting, frequency, sender, family, zone)
//	appNotificationSettings = append(appNotificationSettings, appNotificationSetting)
//}
//return appNotificationSettings, nil

//notificationTypeMap := make(map[string]*entity.EventType)
//for _, nt := range notificationTypes {
//	notificationTypeMap[nt.externalID.String()] = nt
//}
//
//settingsTypesMap := make(map[string][]*entity.EventType)
//for _, relation := range notifySettingsNotifyTypesRelations {
//	if nt, exists := notificationTypeMap[relation.NotifyTypeID]; exists {
//		settingsTypesMap[relation.NotifySettingID] = append(settingsTypesMap[relation.NotifySettingID], nt)
//	}
//}

//appNotificationSettings := make([]*NotificationSettingDetails, len(notificationSettingsIDs))
//appNotificationFrequencies := make([]*ApplicationNotificationFrequency, len(notificationFrequencies))
//appNotificationSenders := make([]*ApplicationUser, len(notificationSenders))
//for k, id := range notificationSettingsIDs {
//	notificationTypesForSetting := settingsTypesMap[id]
//	appNotificationSettings[k] = NewApplicationNotificationSetting(id, notificationTypesForSetting)
//	appNotificationFrequencies[k] = NewApplicationNotificationFrequency(notificationFrequencies[k])
//	appNotificationSenders[k] = NewApplicationUser(notificationSenders[k])
//}

//}

type ApplicationNotificationsSendersByReceiver struct {
	NotificationSenders []Sender
}

type Sender struct {
	Name  string
	Phone string
}

func NewApplicationNotificationsSendersByReceiver(
	senders map[string]string) *ApplicationNotificationsSendersByReceiver {
	notificationSenders := make([]Sender, 0, len(senders))
	for phone, name := range senders {
		notificationSenders = append(notificationSenders, Sender{
			Name:  name,
			Phone: phone,
		})
	}

	return &ApplicationNotificationsSendersByReceiver{
		NotificationSenders: notificationSenders,
	}
}
