package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
)

type GetInvitationResponse struct {
	InvitationID uuid.UUID           `json:"invitation_id"`
	Author       GetUserResponse     `json:"author"`
	Family       GetFamilyResponse   `json:"family"`
	Role         GetRoleByIDResponse `json:"role"`
} // @name InvitationResponse

func NewGetInvitationByReceiverResponse(
	invitation *app_model.ApplicationInvitationToFamilyType) *GetInvitationResponse {
	return &GetInvitationResponse{
		InvitationID: invitation.ID,
		Author:       *NewGetUserResponse(&invitation.Author),
		Family:       *NewGetFamilyResponse(&invitation.Family),
		Role:         *NewGetRoleByIDResponse(&invitation.Role),
	}
}

func NewGetInvitationsByReceiverResponse(
	invitations []*app_model.ApplicationInvitationToFamilyType) []*GetInvitationResponse {
	respInvitations := make([]*GetInvitationResponse, len(invitations))
	for i, invitation := range invitations {
		respInvitations[i] = NewGetInvitationByReceiverResponse(invitation)
	}
	return respInvitations
}

type AddInvitationResponse struct {
	ID uuid.UUID `json:"id"`
} // @name AddInvitationResponse

func NewAddInvitationResponse(id value_object.ID) *AddInvitationResponse {
	return &AddInvitationResponse{ID: id.ToRaw()}
}

//type GetInvitationsByReceiverResponse struct {
//	Invitations []GetInvitationByReceiverResponse `json:"invitations"`
//}
//
//type GetInvitationByReceiverResponse struct {
//	SenderName string `json:"full_name"`
//	Phone      string `json:"phone"`
//	Family     string `json:"family"`
//	Role       string `json:"role"`
//}
//
