package response

//type GetMessengerTypeResponse struct {
//	ID   string `json:"id"`
//	Type string `json:"mes_type"`
//}
//
//func NewGetMessengerTypeResponse(
//	messengerType *app_model.ApplicationMessengerType) *GetMessengerTypeResponse {
//	return &GetMessengerTypeResponse{
//		ID:   messengerType.ID,
//		Type: messengerType.Type,
//	}
//}
//
//func NewGetMessengerTypesResponse(
//	messengerTypes []*app_model.ApplicationMessengerType) []*GetMessengerTypeResponse {
//	var responseTypes []*GetMessengerTypeResponse
//	for _, mt := range messengerTypes {
//		responseType := NewGetMessengerTypeResponse(mt)
//		if responseType != nil {
//			responseTypes = append(responseTypes, responseType)
//		}
//	}
//	return responseTypes
//}
