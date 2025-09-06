package response

import (
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

//type InvitationToFamilyDB struct {
//	ID        uuid.UUID `db:"id"`
//	AuthorID  uuid.UUID `db:"author_id"`
//	RoleID    int       `db:"role_id"`
//	FamilyID  uuid.UUID `db:"family_id"`
//	CreatedAt time.Time `db:"created_at"`
//}

type InvitationToFamilyDB struct {
	ID                  string     `db:"id"`
	InvitationCreatedAt time.Time  `db:"invitation_created_at"`
	InvitationEndedAt   *time.Time `db:"invitation_ended_at"`
	AuthorID            string     `db:"author_id"`
	AuthorFirstName     string     `db:"author_first_name"`
	AuthorLastName      string     `db:"author_last_name"`
	AuthorEmail         string     `db:"email"`
	AuthorLogin         string     `db:"username"`
	AuthorTracking      bool       `db:"tracking"`
	AuthorCreatedAt     time.Time  `db:"author_created_at"`
	AuthorUpdatedAt     time.Time  `db:"author_updated_at"`
	TelegramID          *string    `db:"telegram_id"`
	FamilyID            string     `db:"family_id"`
	FamilyName          string     `db:"family_name"`
	FamilyCreatedAt     time.Time  `db:"family_created_at"`
	FamilyUpdatedAt     time.Time  `db:"family_updated_at"`
	FamilyEndedAt       *time.Time `db:"family_ended_at"`
	RoleID              int        `db:"role_id"`
	RoleName            string     `db:"role_name"`
}

type ShortInvitationToFamilyDB struct {
	ID        string    `db:"id"`
	AuthorID  string    `db:"author_id"`
	RoleID    int       `db:"role_id"`
	FamilyID  string    `db:"family_id"`
	CreatedAt time.Time `db:"created_at"`
}

func InvitationDbToEntity(db *InvitationToFamilyDB) (*aggregate.InvitationToFamily, error) {
	id, err := value_object.NewIDFromString(db.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to id")
	}
	userID, err := value_object.NewIDFromString(db.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user id")
	}
	familyID, err := value_object.NewIDFromString(db.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to family id")
	}

	telegramExists := false
	if db.TelegramID != nil {
		telegramExists = true
	}

	return &aggregate.InvitationToFamily{
		ID: id,
		Author: &entity.User{
			ID:           userID,
			Name:         value_object.FullName{FirstName: value_object.FirstName(db.AuthorFirstName), LastName: value_object.LastName(db.AuthorLastName)},
			Email:        value_object.Email(db.AuthorEmail),
			Login:        value_object.Login(db.AuthorLogin),
			ExternalID:   value_object.NewID(),
			Tracking:     db.AuthorTracking,
			CreatedAt:    db.AuthorCreatedAt,
			UpdatedAt:    db.AuthorUpdatedAt,
			EndedAt:      nil,
			TgIntegrated: telegramExists,
		},
		Role: &entity.Role{
			ID:   db.RoleID,
			Name: value_object.RoleName(db.RoleName),
		},
		Family: &entity.Family{
			ID:        familyID,
			Name:      value_object.FamilyName(db.FamilyName),
			CreatedAt: db.FamilyCreatedAt,
			UpdatedAt: db.FamilyUpdatedAt,
			EndedAt:   db.FamilyEndedAt,
		},
		CreatedAt: db.InvitationCreatedAt,
		EndedAt:   db.InvitationEndedAt,
	}, nil
}

func InvitationDbListToEntityList(dbList []*InvitationToFamilyDB) ([]*aggregate.InvitationToFamily, error) {
	var invitations []*aggregate.InvitationToFamily
	for _, db := range dbList {
		invitation, err := InvitationDbToEntity(db)
		if err != nil {
			return nil, fmt.Errorf("failed to convert db data to entity: %w", err)
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}
