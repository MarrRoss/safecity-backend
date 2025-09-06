package db

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type FamilyMembershipRepository interface {
	AddFamilyMembership(ctx context.Context, familyMembership *aggregate.FamilyMembership) error
	GetFamilyMembershipByID(ctx context.Context, id value_object.ID) (*response.FamilyMembershipDB, error)
	GetFamiliesByUserID(ctx context.Context, userID value_object.ID) ([]*entity.Family, error)
	MembershipExists(ctx context.Context, userID, familyID value_object.ID) (bool, error)
	//GetMembershipsByUserID(ctx context.Context, userID value_object.ID) ([]*response.FamilyMembershipDB, error)
	CheckUsersBelongToFamilyByZone(ctx context.Context, zoneID, senderID, receiverID value_object.ID) (bool, error)
	CheckUsersShareCommonActiveFamily(ctx context.Context, senderID, receiverID value_object.ID) (bool, error)
	GetMembershipsByFamilyID(ctx context.Context, familyID value_object.ID, usersIDs *[]string,
	) ([]*response.FamilyMembershipDB, error)
	GetMembershipsParticipantsByUser(ctx context.Context,
		userID value_object.ID) ([]*response.FamilyMembershipParticipants, error)
	GetAvailableNotificationSenders(
		ctx context.Context, receiverID, familyID value_object.ID) ([]*entity.User, error)
	DeleteFamilyMembership(ctx context.Context, membershipID value_object.ID) error
	UpdateFamilyMembershipRole(ctx context.Context, membershipID value_object.ID, roleID int) error
	//GetMembershipByUserAndZone(ctx context.Context, userID, zoneID value_object.externalID) (*response.FamilyMembershipDB, error)
	//DeleteFamilyMembership(ctx context.Context, id value_object.externalID) error
	//DeleteFamilyMembershipWithCheck(ctx context.Context, membershipID value_object.externalID) error
}
