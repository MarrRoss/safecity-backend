package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
)

type GetNotificationSettingResponse struct {
	ID         uuid.UUID                         `json:"notify_setting_id"`
	Frequency  *GetNotificationFrequencyResponse `json:"frequency,omitempty"`
	Sender     GetUserResponse                   `json:"sender"`
	MinBattery *int                              `json:"min_battery,omitempty"`
	//NotificationTypes []*ApplicationNotificationType
	//MessengerTypes    []*ApplicationMessengerType
} // @name GetNotificationSettingResponse

func NewGetNotificationSettingResponse(setting *app_model.ApplicationNotificationSettings) *GetNotificationSettingResponse {
	return &GetNotificationSettingResponse{
		ID:         setting.ID,
		Frequency:  NewGetNotificationFrequencyResponse(setting.Frequency),
		Sender:     *NewGetUserResponse(setting.Sender),
		MinBattery: setting.MinBattery,
	}
}

func NewGetNotificationSettingsResponse(settings []*app_model.ApplicationNotificationSettings) []*GetNotificationSettingResponse {
	respSettings := make([]*GetNotificationSettingResponse, len(settings))
	for i, setting := range settings {
		respSettings[i] = NewGetNotificationSettingResponse(setting)
	}
	return respSettings
}

//type GetNotificationSettingResponse struct {
//	ID                string                           `json:"id"`
//	Frequency         GetNotificationFrequencyResponse `json:"frequency"`
//	Sender            GetUserResponse                  `json:"sender"`
//	NotificationTypes []*GetNotificationTypeResponse   `json:"notify_types"`
//	//MessengerTypes    []*GetMessengerTypeResponse      `json:"mes_types"`
//	Family *GetFamilyResponse `json:"family"`
//	Zone   *GetZoneResponse   `json:"zone"`
//}

type GetNotificationsSendersByReceiverResponse struct {
	NotificationSenders []Sender `json:"notify_senders"`
}

type Sender struct {
	Name  string `json:"full_name"`
	Phone string `json:"phone"`
}

//func NewGetNotificationSettingResponse(
//	notificationSetting *app_model.NotificationSettingDetails,
//) *GetNotificationSettingResponse {
//	return &GetNotificationSettingResponse{
//		ID:        notificationSetting.ID.String(),
//		Frequency: *NewGetNotificationFrequencyResponse(notificationSetting.Frequency),
//		Sender:    *NewGetUserResponse(notificationSetting.Sender),
//		//NotificationTypes: NewGetNotificationTypesResponse(
//		//	notificationSetting.NotificationTypes),
//		//MessengerTypes: NewGetMessengerTypesResponse(
//		//	notificationSetting.MessengerTypes),
//		//Family: NewGetFamilyResponse(notificationSetting.Family),
//		//Zone:   NewGetZoneResponse(notificationSetting.Zone),
//	}
//}
//
//func NewGetNotificationSettingsResponse(
//	notificationSettings []*app_model.NotificationSettingDetails,
//) []*GetNotificationSettingResponse {
//	responseSettings := make([]*GetNotificationSettingResponse, 0, len(notificationSettings))
//	for _, setting := range notificationSettings {
//		responseSettings = append(responseSettings, NewGetNotificationSettingResponse(setting))
//	}
//
//	return responseSettings
//}

func NewGetNotificationsSendersByReceiverResponse(
	senders *app_model.ApplicationNotificationsSendersByReceiver) *GetNotificationsSendersByReceiverResponse {
	var respSenders []Sender
	for _, sender := range senders.NotificationSenders {
		respSenders = append(respSenders, Sender{
			Name:  sender.Name,
			Phone: sender.Phone,
		})
	}
	return &GetNotificationsSendersByReceiverResponse{
		NotificationSenders: respSenders,
	}
}

type AddNotificationSettingResponse struct {
	ID uuid.UUID `json:"id"`
} // @name AddNotificationSettingResponse

func NewAddNotificationSettingResponse(id value_object.ID) *AddNotificationSettingResponse {
	return &AddNotificationSettingResponse{ID: id.ToRaw()}
}
