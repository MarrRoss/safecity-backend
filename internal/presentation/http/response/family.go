package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
)

type GetFamilyResponse struct {
	ID   uuid.UUID `json:"family_id"`
	Name string    `json:"family_name"`
}

func NewGetFamilyResponse(family *app_model.ApplicationFamily) *GetFamilyResponse {
	return &GetFamilyResponse{
		ID:   family.ID,
		Name: family.Name,
	}
}

func NewGetFamiliesResponse(families []*app_model.ApplicationFamily) []*GetFamilyResponse {
	respFamilies := make([]*GetFamilyResponse, len(families))
	for i, family := range families {
		respFamilies[i] = NewGetFamilyResponse(family)
	}
	return respFamilies
}

type AddFamilyResponse struct {
	FamilyID     uuid.UUID `json:"family_id"`
	MembershipID uuid.UUID `json:"membership_id"`
} // @name AddFamilyResponse

func NewAddFamilyResponse(familyID, membershipID value_object.ID) *AddFamilyResponse {
	return &AddFamilyResponse{
		FamilyID:     familyID.ToRaw(),
		MembershipID: membershipID.ToRaw(),
	}
}

type AddFamilyZoneResponse struct {
	ID uuid.UUID `json:"id"`
} // @name AddFamilyZoneResponse

func NewAddFamilyZoneResponse(id value_object.ID) *AddFamilyZoneResponse {
	return &AddFamilyZoneResponse{ID: id.ToRaw()}
}
