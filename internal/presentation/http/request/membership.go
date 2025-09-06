package request

import "github.com/google/uuid"

//type GetUserMembershipsRequest struct {
//	ID uuid.UUID `params:"id" format:"uuid"`
//} // @name GetUserMembershipsRequest

type GetMembershipsByFamilyRequest struct {
	FamilyID string `json:"family_id" path:"family_id"`
}

type GetInvitationByIDRequest struct {
	ID uuid.UUID `json:"id" format:"uuid"`
} // @name GetInvitationByIDRequest

type GetFamilyMembershipByIDRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name GetFamilyMembershipByIDRequest

type RoleType int // @name RoleType
const (
	RoleParentType RoleType = 1
	RoleChildType  RoleType = 2
)

type ChangeFamilyMembershipRequest struct {
	RoleID RoleType `json:"role_id" enum:"1,2"`
} // @name ChangeFamilyMembershipRequest

//type AddMembershipRequest struct {
//	UserID   uuid.UUID `json:"user_id" format:"uuid"`
//	RoleID   int       `json:"role_id" format:"int"`
//	FamilyID uuid.UUID `json:"family_id" format:"uuid"`
//}

type AddMembershipByInvitationRequest struct {
	ID string `json:"invitation_id" path:"invitation_id"`
}

type DeleteUserMembershipRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name DeleteUserMembershipRequest
