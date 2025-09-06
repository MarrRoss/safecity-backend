package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
)

type AddUserLocationResponse struct {
	LocationLogID uuid.UUID   `json:"location_log_id"`
	NotifyLogID   []uuid.UUID `json:"notify_log_id"`
} // @AddUserLocationResponse

func NewAddUserLocationResponse(locationLogID value_object.ID, notifyLogsIDs []value_object.ID) *AddUserLocationResponse {
	respNotifyLogsIDs := make([]uuid.UUID, len(notifyLogsIDs))
	for i, notifyLogID := range notifyLogsIDs {
		respNotifyLogsIDs[i] = notifyLogID.ToRaw()
	}
	return &AddUserLocationResponse{
		LocationLogID: locationLogID.ToRaw(),
		NotifyLogID:   respNotifyLogsIDs,
	}
}

type GetUserLocationResponse struct {
	User     GetUserResponse     `json:"user"`
	Location GetLocationResponse `json:"location"`
}

type GetLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewGetLocationResponse(location *app_model.ApplicationLocation) *GetLocationResponse {
	return &GetLocationResponse{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}

func NewGetUserLocationResponse(location *app_model.ApplicationUserLocation) *GetUserLocationResponse {
	return &GetUserLocationResponse{
		User:     *NewGetUserResponse(location.User),
		Location: *NewGetLocationResponse(location.Location),
	}
}

func NewGetUserLocationResponses(locations []*app_model.ApplicationUserLocation) []*GetUserLocationResponse {
	appLocations := make([]*GetUserLocationResponse, len(locations))
	for i, location := range locations {
		appLocations[i] = NewGetUserLocationResponse(location)
	}
	return appLocations
}

//type AddFamilyResponse struct {
//	FamilyID     uuid.UUID `json:"family_id"`
//	MembershipID uuid.UUID `json:"membership_id"`
//} // @name AddFamilyResponse
//
//func NewAddFamilyResponse(familyID, membershipID value_object.ID) *AddFamilyResponse {
//	return &AddFamilyResponse{
//		FamilyID:     familyID.ToRaw(),
//		MembershipID: membershipID.ToRaw(),
//	}
//}
