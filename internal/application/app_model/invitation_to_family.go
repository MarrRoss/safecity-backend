package app_model

import (
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"github.com/google/uuid"
)

type ApplicationInvitationToFamilyType struct {
	ID     uuid.UUID
	Author ApplicationUser
	Family ApplicationFamily
	Role   ApplicationRole
}

func NewApplicationInvitationToFamilyType(invitation *aggregate.InvitationToFamily) *ApplicationInvitationToFamilyType {
	return &ApplicationInvitationToFamilyType{
		ID:     invitation.ID.ToRaw(),
		Author: *NewApplicationUser(invitation.Author),
		Family: *NewApplicationFamily(invitation.Family),
		Role:   *NewApplicationRole(invitation.Role),
	}
}

func NewApplicationInvitationToFamilyTypes(invitations []*aggregate.InvitationToFamily) []*ApplicationInvitationToFamilyType {
	appInvitations := make([]*ApplicationInvitationToFamilyType, len(invitations))
	for i, invitation := range invitations {
		appInvitations[i] = NewApplicationInvitationToFamilyType(invitation)
	}
	return appInvitations
}

type ApplicationInvitationToFamily struct {
	ID uuid.UUID
}
