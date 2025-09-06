package request

import "github.com/google/uuid"

type GetInvitationsToFamilyByReceiverIdRequest struct {
	UserID string `json:"user_id" path:"user_id"`
}

type AddInvitationToFamilyRequest struct {
	RoleID   int       `json:"role_id" format:"int"`
	FamilyID uuid.UUID `json:"family_id" format:"uuid"`
} // @name AddInvitationToFamilyRequest

type AddMembershipCommand struct {
	UserID uuid.UUID `json:"user_id" format:"uuid"`
} // @name AddMembershipCommand
