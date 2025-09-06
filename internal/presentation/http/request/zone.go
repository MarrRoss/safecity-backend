package request

import "github.com/google/uuid"

type GetZoneByIDRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name GetZoneByIDRequest

type GetZonesByUserIDRequest struct {
	UserID uuid.UUID `params:"id" format:"uuid"`
} // @name GetZonesByUserIDRequest

type AddZoneRequest struct {
	Name       string          `json:"name"`
	Boundaries []LocationModel `json:"boundaries"`
	Safety     bool            `json:"safety"`
} // @name AddZoneRequest

type LocationModel struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ChangeZoneRequest struct {
	Name       *string          `json:"name"`
	Boundaries *[]LocationModel `json:"boundaries"`
	Safety     *bool            `json:"safety"`
} // @name ChangeZoneRequest

type DeleteZoneRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name DeleteZoneRequest
