package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
	"time"
)

type GetUserResponse struct {
	ID           uuid.UUID `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Login        string    `json:"login"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Tracking     bool      `json:"tracking"`
	TgIntegrated bool      `json:"tg_integrated"`
} // @name GetUserResponse

func NewGetUserResponse(user *app_model.ApplicationUser) *GetUserResponse {
	return &GetUserResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Login:        user.Login,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Tracking:     user.Tracking,
		TgIntegrated: user.TgIntegrated,
	}
}

func NewGetUsersResponse(users []*app_model.ApplicationUser) []*GetUserResponse {
	respUsers := make([]*GetUserResponse, len(users))
	for i, user := range users {
		respUsers[i] = NewGetUserResponse(user)
	}
	return respUsers
}

type AddUserResponse struct {
	ID uuid.UUID `json:"id"`
} // @name AddUserResponse

func NewAddUserResponse(id value_object.ID) *AddUserResponse {
	return &AddUserResponse{ID: id.ToRaw()}
}

//type GetUserLocationsResponse struct {
//	Locations []*LocationResponse `json:"locations"`
//}

//type LocationResponse struct {
//	Latitude  float64 `json:"latitude"`
//	Longitude float64 `json:"longitude"`
//	Time      string  `json:"time"`
//}

//func NewGetUserLocationsResponse(locations []*app_model.ApplicationLocationTime) *GetUserLocationsResponse {
//	respLocations := make([]*LocationResponse, len(locations))
//	for i, loc := range locations {
//		respLocations[i] = &LocationResponse{
//			Latitude:  loc.Latitude,
//			Longitude: loc.Longitude,
//			Time:      loc.Time,
//		}
//	}
//	return &GetUserLocationsResponse{Locations: respLocations}
//}
