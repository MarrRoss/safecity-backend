package db

import (
	response2 "awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type InvitationToFamilyRepository interface {
	AddInvitationToFamily(ctx context.Context, invitation *aggregate.InvitationToFamily) error
	GetInvitationToFamily(ctx context.Context, invitationID value_object.ID) (*response2.ShortInvitationToFamilyDB, error)
	//GetInvitationsToFamilyByReceiverID(ctx context.Context,
	//	receiverID value_object.ID) ([]*response.InvitationToFamilyDB, error)
	GetInvitationToFamilyByID(ctx context.Context, invitationID value_object.ID) (*aggregate.InvitationToFamily, error)
	GetPendingInvitationsByUser(ctx context.Context, userID value_object.ID) ([]*aggregate.InvitationToFamily, error)
	AddInviteActivation(ctx context.Context, invitationID, userID value_object.ID) error
	//AcceptInvitationToFamily(ctx context.Context, invitationID value_object.ID) error
}
