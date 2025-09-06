package app_model

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"github.com/google/uuid"
)

//type ApplicationFamilyMemberships struct {
//	Memberships []*ApplicationFamilyMembership
//}

type ApplicationFamilyMembership struct {
	UserMembershipID uuid.UUID
	Family           *ApplicationFamily
	UserRole         string
	Participants     []*ApplicationMembershipParticipant
}

type ApplicationMembershipParticipant struct {
	MembershipID uuid.UUID
	User         *ApplicationUser
	RoleName     string
}

func NewApplicationFamilyMemberships(
	currentUserID string,
	memberships []*response.FamilyMembershipParticipants,
	users []*entity.User,
	families []*entity.Family,
) ([]*ApplicationFamilyMembership, error) {
	// Создаем мапу для быстрого поиска пользователей по ID
	userMap := make(map[string]*entity.User)
	for _, user := range users {
		userMap[user.ID.String()] = user
	}
	// Создаем мапу для быстрого поиска семей по ID
	familyMap := make(map[string]*entity.Family)
	for _, fam := range families {
		familyMap[fam.ID.String()] = fam
	}

	// Группируем членства по familyID
	grouped := make(map[string]*ApplicationFamilyMembership)
	for _, mem := range memberships {
		u, ok := userMap[mem.UserID]
		if !ok {
			return nil, fmt.Errorf("user not found for membership, id %s", mem.UserID)
		}
		f, ok := familyMap[mem.FamilyID]
		if !ok {
			return nil, fmt.Errorf("family not found for membership, id %s", mem.FamilyID)
		}
		// Если для данной семьи еще не создан объект, создаем новый
		if _, exists := grouped[mem.FamilyID]; !exists {
			grouped[mem.FamilyID] = &ApplicationFamilyMembership{
				Family:       NewApplicationFamily(f),
				Participants: []*ApplicationMembershipParticipant{},
			}
		}
		// Если найдено членство текущего пользователя – сохраняем его роль и membershipID
		if mem.UserID == currentUserID {
			membershipID, err := value_object.NewIDFromString(mem.MembershipID)
			if err != nil {
				return nil, fmt.Errorf("invalid membershipID: %w", err)
			}
			grouped[mem.FamilyID].UserRole = mem.RoleName
			grouped[mem.FamilyID].UserMembershipID = membershipID.ToRaw()
		} else {
			// Для остальных участников создаем запись участника
			membershipID, err := value_object.NewIDFromString(mem.MembershipID)
			if err != nil {
				return nil, fmt.Errorf("invalid membershipID: %w", err)
			}
			participant := &ApplicationMembershipParticipant{
				MembershipID: membershipID.ToRaw(),
				User:         NewApplicationUser(u),
				RoleName:     mem.RoleName,
			}
			grouped[mem.FamilyID].Participants = append(grouped[mem.FamilyID].Participants, participant)
		}
	}

	// Преобразуем сгруппированные данные в срез
	var result []*ApplicationFamilyMembership
	for _, v := range grouped {
		result = append(result, v)
	}
	return result, nil
}

//type ApplicationFamilyMembership struct {
//	externalID       uuid.UUID
//	User     *ApplicationUser
//	RoleName string
//	Family   *ApplicationFamily
//}

//func NewApplicationFamilyMembership(
//	membership *response.FamilyMembershipParticipants,
//	user *entity.User,
//	family *entity.Family) (*ApplicationFamilyMembership, error) {
//	id, err := value_object.NewIDFromString(membership.MembershipID)
//	if err != nil {
//		return nil, domain.ErrInvalidID
//	}
//	return &ApplicationFamilyMembership{
//		externalID:       id.ToRaw(),
//		User:     NewApplicationUser(user),
//		RoleName: membership.RoleName,
//		Family:   NewApplicationFamily(family),
//	}, nil
//}

//func NewApplicationFamilyMemberships(
//	memberships []*response.FamilyMembershipParticipants,
//	users []*entity.User,
//	families []*entity.Family) (
//	[]*ApplicationFamilyMembership,
//	error) {
//	userMap := make(map[string]*entity.User)
//	for _, user := range users {
//		userMap[user.externalID.String()] = user
//	}
//	familyMap := make(map[string]*entity.Family)
//	for _, family := range families {
//		familyMap[family.externalID.String()] = family
//	}
//
//	var applicationFamilyMemberships []*ApplicationFamilyMembership
//
//	for _, membership := range memberships {
//		user, userExists := userMap[membership.UserID]
//		if !userExists {
//			return nil, fmt.Errorf("externalID %s: %w", membership.UserID, application.ErrMembershipUserNotFound)
//		}
//
//		family, familyExists := familyMap[membership.FamilyID]
//		if !familyExists {
//			return nil, fmt.Errorf("externalID %s: %w", membership.FamilyID, application.ErrMembershipFamilyNotFound)
//		}
//
//		applicationMembership, err := NewApplicationFamilyMembership(membership, user, family)
//		if err != nil {
//			return nil, err
//		}
//		applicationFamilyMemberships = append(applicationFamilyMemberships, applicationMembership)
//	}
//
//	return applicationFamilyMemberships, nil

//appMemberships := make([]*ApplicationFamilyMembership, len(memberships))
//appUsers := make([]*ApplicationUser, len(users))
//appFamilies := make([]*ApplicationFamily, len(families))
//for i, membership := range memberships {
//	appMemberships[i] = NewApplicationFamilyMembership(membership)
//	appUsers[i] = NewApplicationUser(users[i])
//	appFamilies[i] = NewApplicationFamily(families[i])
//}
//if len(appMemberships) != len(appUsers) || len(appMemberships) != len(appFamilies) {
//	return nil, nil, nil, errors.New("memberships, users and families slices must have the same length")
//}
//return appMemberships, appUsers, appFamilies, nil
//}
