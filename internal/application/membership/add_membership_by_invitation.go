package membership

import (
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type AddMembershipByInvitationCommand struct {
	ID string
}

type AddMembershipByInvitationHandler struct {
	invitationStorage db.InvitationToFamilyRepository
	membershipStorage db.FamilyMembershipRepository
	userStorage       db.UserRepository
	roleStorage       db.RoleRepository
	familyStorage     db.FamilyRepository
}

func NewAddMembershipByInvitationHandler(
	invitationStorage db.InvitationToFamilyRepository,
	membershipStorage db.FamilyMembershipRepository,
	userStorage db.UserRepository,
	roleStorage db.RoleRepository,
	familyStorage db.FamilyRepository) *AddMembershipByInvitationHandler {
	return &AddMembershipByInvitationHandler{
		invitationStorage: invitationStorage,
		membershipStorage: membershipStorage,
		userStorage:       userStorage,
		roleStorage:       roleStorage,
		familyStorage:     familyStorage,
	}
}

func (h *AddMembershipByInvitationHandler) Handle(ctx context.Context,
	cmd AddMembershipByInvitationCommand) (string, error) {
	//invitationID, err := value_object.NewIDFromString(cmd.ID)
	//if err != nil {
	//	return "", fmt.Errorf("invalid invitation id: %w", err)
	//}
	//invitation, err := h.invitationStorage.GetInvitationToFamily(ctx, invitationID)
	//if err != nil {
	//	return "", fmt.Errorf("failed to get invitation: %w", err)
	//}
	//
	//familyID, err := value_object.NewIDFromString(invitation.FamilyID)
	//if err != nil {
	//	return "", fmt.Errorf("invalid family id: %w", err)
	//}
	//family, err := h.familyStorage.GetFamily(ctx, familyID)
	//if err != nil {
	//	return "", fmt.Errorf("failed to get family: %w", err)
	//}
	//
	//exists, err := h.membershipStorage.MembershipExists(ctx, receiver.ID, familyID)
	//if err != nil {
	//	return "", fmt.Errorf("failed to check membership exists: %w", err)
	//}
	//if exists {
	//	return "", errors.New("this membership already exists in system")
	//}
	//
	//roleID, err := strconv.Atoi(invitation.RoleID)
	//if err != nil {
	//	return "", fmt.Errorf("failed to convert role id to int: %w", err)
	//}
	//role, err := h.roleStorage.GetRole(ctx, roleID)
	//if err != nil {
	//	return "", fmt.Errorf("failed to get role: %w", err)
	//}
	//membership, err := aggregate.NewFamilyMembership(receiver, role, family)
	//if err != nil {
	//	return "", fmt.Errorf("failed to create family membership: %w", err)
	//}
	//err = h.membershipStorage.AddFamilyMembership(ctx, membership)
	//if err != nil {
	//	return "", fmt.Errorf("failed to add family membership to storage: %w", err)
	//}
	//
	//err = h.invitationStorage.AcceptInvitationToFamily(ctx, invitationID)
	//if err != nil {
	//	return "", fmt.Errorf("failed to accept invitation: %w", err)
	//}

	return "", nil
}
