package response

type GetNotificationTypeResponse struct {
	ID         string `json:"id"`
	NotifyType string `json:"notify_type"`
}

//type GetNotificationTypesResponse struct {
//	EventType []GetNotificationTypeResponse `json:"notify_types"`
//}

//type GetNotificationTypeResponse struct {
//	NotificationType string `json:"notify_type"`
//}

//func NewGetNotificationTypeResponse(
//	notificationType *app_model.ApplicationNotificationType,
//) *GetNotificationTypeResponse {
//	if notificationType == nil {
//		return nil
//	}
//	return &GetNotificationTypeResponse{
//		ID:         notificationType.ID,
//		NotifyType: notificationType.Type,
//	}
//}

//func NewGetNotificationTypesResponse(
//	notificationTypes []*app_model.ApplicationNotificationType) *GetNotificationTypesResponse {
//	respNotTypes := make([]GetNotificationTypeResponse, len(notificationTypes))
//	for i, notType := range notificationTypes {
//		respNotTypes[i] = *NewGetNotificationTypeResponse(notType)
//	}
//	return &GetNotificationTypesResponse{
//		EventType: respNotTypes,
//	}
//}

//func NewGetNotificationTypesResponse(
//	notificationTypes []*app_model.ApplicationNotificationType,
//) []*GetNotificationTypeResponse {
//	var responseTypes []*GetNotificationTypeResponse
//	for _, nt := range notificationTypes {
//		responseType := NewGetNotificationTypeResponse(nt)
//		if responseType != nil {
//			responseTypes = append(responseTypes, responseType)
//		}
//	}
//	return responseTypes
//}
