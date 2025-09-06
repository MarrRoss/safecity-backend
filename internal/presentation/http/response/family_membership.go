package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
)

//type GetFamilyMembershipsResponse struct {
//	Memberships []*GetFamilyMembershipResponse `json:"memberships"`
//}

//type GetFamilyMembershipsSuccessResponse struct {
//	BaseStatusErrorResponse
//	Data []*GetFamilyMembershipResponse `json:"data"`
//}

func ConvertRoleName(v string) string {
	switch v {
	case "adult":
		return "Родитель"
	case "child":
		return "Ребёнок"
	}
	return v
}

type GetMembershipParticipantResponse struct {
	MembershipID uuid.UUID       `json:"membership_id"`
	User         GetUserResponse `json:"user"`
	RoleName     string          `json:"role_name"`
}
type GetFamilyMembershipResponse struct {
	MembershipID uuid.UUID                           `json:"membership_id"`
	Family       GetFamilyResponse                   `json:"family"`
	MyRole       string                              `json:"my_role"`
	Participants []*GetMembershipParticipantResponse `json:"participants"`
} // @name GetFamilyMembershipResponse

func NewGetMembershipParticipantResponse(participant *app_model.ApplicationMembershipParticipant) *GetMembershipParticipantResponse {
	return &GetMembershipParticipantResponse{
		MembershipID: participant.MembershipID,
		User:         *NewGetUserResponse(participant.User),
		RoleName:     ConvertRoleName(participant.RoleName),
	}
}

func NewGetFamilyMembershipResponse(membership *app_model.ApplicationFamilyMembership) *GetFamilyMembershipResponse {
	var participantsResp []*GetMembershipParticipantResponse
	for _, p := range membership.Participants {
		participantsResp = append(participantsResp, NewGetMembershipParticipantResponse(p))
	}
	var participants []*GetMembershipParticipantResponse
	if len(participantsResp) > 0 {
		participants = participantsResp
	} else {
		participants = []*GetMembershipParticipantResponse{}
	}

	return &GetFamilyMembershipResponse{
		MembershipID: membership.UserMembershipID,
		Family:       *NewGetFamilyResponse(membership.Family),
		MyRole:       ConvertRoleName(membership.UserRole),
		Participants: participants,
	}
}

func NewGetFamilyMembershipsResponse(
	memberships []*app_model.ApplicationFamilyMembership,
) []*GetFamilyMembershipResponse {
	var respMemberships []*GetFamilyMembershipResponse
	for _, membership := range memberships {
		respMemberships = append(respMemberships, NewGetFamilyMembershipResponse(membership))
	}
	if len(respMemberships) == 0 {
		return []*GetFamilyMembershipResponse{}
	}
	return respMemberships
}

//type GetFamilyMembershipResponse struct {
//	externalID       uuid.UUID         `json:"membership_id"`
//	User     GetUserResponse   `json:"user"`
//	RoleName string            `json:"role_name"`
//	Family   GetFamilyResponse `json:"family"`
//} // @name GetFamilyMembershipResponse

//func NewGetFamilyMembershipResponse(membership *app_model.ApplicationFamilyMembership) *GetFamilyMembershipResponse {
//	return &GetFamilyMembershipResponse{
//		externalID:       membership.externalID,
//		User:     *NewGetUserResponse(membership.User),
//		RoleName: membership.RoleName,
//		Family:   *NewGetFamilyResponse(membership.Family),
//	}
//}

//func NewGetFamilyMembershipsResponse(
//	memberships []*app_model.ApplicationFamilyMembership) []*GetFamilyMembershipResponse {
//	var respMemberships []*GetFamilyMembershipResponse
//	for _, membership := range memberships {
//		respMemberships = append(respMemberships, NewGetFamilyMembershipResponse(membership))
//	}
//	return respMemberships
//}

type AddFamilyMembershipResponse struct {
	ID uuid.UUID `json:"id"`
} // @name AddFamilyMembershipResponse

func NewAddFamilyMembershipResponse(id value_object.ID) *AddFamilyMembershipResponse {
	return &AddFamilyMembershipResponse{ID: id.ToRaw()}
}
