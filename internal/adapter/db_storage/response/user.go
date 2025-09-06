package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type UserDB struct {
	UserId     string    `db:"id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	Email      string    `db:"email"`
	Login      string    `db:"username"`
	Tracking   bool      `db:"tracking"`
	TelegramID *string   `db:"telegram_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func UserDbToEntity(db *UserDB) (*entity.User, error) {
	userID, err := value_object.NewIDFromString(db.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user id")
	}
	userFirstName, err := value_object.NewFirstName(db.FirstName)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user first name")
	}
	userLastName, err := value_object.NewLastName(db.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user last name")
	}
	userName := value_object.NewFullName(userFirstName, userLastName)
	userEmail, err := value_object.NewEmail(db.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user email")
	}
	userLogin, err := value_object.NewLogin(db.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user login")
	}

	telegramExists := false
	if db.TelegramID != nil {
		telegramExists = true
	}

	return &entity.User{
		ID:           userID,
		Name:         userName,
		Email:        userEmail,
		Login:        userLogin,
		Tracking:     db.Tracking,
		TgIntegrated: telegramExists,
		CreatedAt:    db.CreatedAt,
		UpdatedAt:    db.UpdatedAt,
	}, nil
}

func UserDbListToEntityList(dbList []*UserDB) ([]*entity.User, error) {
	var users []*entity.User
	for _, db := range dbList {
		user, err := UserDbToEntity(db)
		if err != nil {
			return nil, fmt.Errorf("failed to convert db data to entity")
		}
		users = append(users, user)
	}
	return users, nil
}
