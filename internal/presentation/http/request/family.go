package request

import "github.com/google/uuid"

type GetFamilyByIDRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name GetFamilyByIDRequest

type GetUserFamiliesRequest struct {
	UserID string `json:"user_id" path:"user_id"`
}

type GetFamilyZonesRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name GetFamilyZonesRequest

type AddFamilyRequest struct {
	Name string `json:"name"`
} // @name AddFamilyRequest

type AddFamilyZoneRequest struct {
	ZoneID uuid.UUID `json:"zone_id" format:"uuid"`
} // @name AddFamilyZoneRequest

type ChangeFamilyRequest struct {
	Name string `json:"name"`
} // @name ChangeFamilyRequest

type DeleteFamilyRequest struct {
	ID string `json:"id" path:"id"`
}

type DeleteFamilyZoneRequest struct {
	FamilyID uuid.UUID `params:"family_id" format:"uuid"`
	ZoneID   uuid.UUID `params:"zone_id" format:"uuid"`
} // @name DeleteFamilyZoneRequest

type DeleteFamilyZonesRequest struct {
	ID string `json:"id" path:"id"`
}
