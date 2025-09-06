package entity

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"time"
)

type User struct {
	ID           value_object.ID
	Name         value_object.FullName
	Email        value_object.Email
	Login        value_object.Login
	ExternalID   value_object.ID
	Tracking     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	EndedAt      *time.Time
	TgIntegrated bool
}

// BirthDate  value_object.BirthDate
// Phone      value_object.Phone

func NewUser(
	name value_object.FullName,
	//bd value_object.BirthDate,
	email value_object.Email,
	//ph value_object.Phone,
	login value_object.Login,
	extID value_object.ID,
) (*User, error) {
	id := value_object.NewID()
	now := time.Now()
	newUser := User{
		ID:           id,
		Name:         name,
		Email:        email,
		Login:        login,
		ExternalID:   extID,
		Tracking:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		EndedAt:      nil,
		TgIntegrated: false,
	}
	return &newUser, nil
}

func (u *User) ChangeFirstName(name value_object.FirstName) error {
	if u.EndedAt != nil {
		return domain.ErrUserIsDeleted
	}
	u.Name.FirstName = name
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ChangeLastName(name value_object.LastName) error {
	if u.EndedAt != nil {
		return domain.ErrUserIsDeleted
	}
	u.Name.LastName = name
	u.UpdatedAt = time.Now()
	return nil
}

//func (u *User) ChangePatronymic(name value_object.Patronymic) error {
//	if u.EndedAt != nil {
//		return domain.ErrUserIsDeleted
//	}
//	u.Name.Patronymic = name
//	u.UpdatedAt = time.Now()
//	return nil
//}

//func (u *User) ChangeBirthDate(bd value_object.BirthDate) error {
//	if u.EndedAt != nil {
//		return domain.ErrUserIsDeleted
//	}
//	u.BirthDate = bd
//	u.UpdatedAt = time.Now()
//	return nil
//}

func (u *User) ChangeTracking(tracking bool) error {
	if u.EndedAt != nil {
		return domain.ErrUserIsDeleted
	}
	u.Tracking = tracking
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) StopExistence() error {
	if u.EndedAt != nil {
		return domain.ErrUserIsDeleted
	}
	u.Tracking = false
	timeNow := time.Now()
	u.EndedAt = &timeNow
	u.UpdatedAt = time.Now()
	return nil
}

//func (u *User) CheckAdult() (bool, error) {
//	if u.EndedAt != nil {
//		return false, domain.ErrUserIsDeleted
//	}
//	birthDate := time.Time(u.BirthDate)
//	now := time.Now()
//	years := now.Year() - birthDate.Year()
//	if years > 18 || (years == 18 && (now.YearDay() >= birthDate.YearDay())) {
//		return true, nil
//	}
//	return false, nil
//}

//func (u *User) RestoreExistence() error {
//	if u.EndedAt == nil {
//		return domain.ErrUserIsDeleted
//	}
//	u.EndedAt = nil
//	u.UpdatedAt = time.Now()
//	return nil
//}

//func (u *User) ChangePassword(pw value_object.Password) error {
//	if u.EndedAt != nil {
//		return domain.ErrUserIsDeleted
//	}
//	u.Password = pw
//	u.UpdatedAt = time.Now()
//	return nil
//}

//func (u *User) ChangeEmail(email value_object.Email) error {
//	if u.EndedAt != nil {
//		return domain.ErrUserIsDeleted
//	}
//	u.Email = &email
//	u.UpdatedAt = time.Now()
//	return nil
//}

//func (u *User) ChangePhone(phone value_object.Phone) error {
//	if u.EndedAt != nil {
//		return errors.New("user is deleted")
//	}
//	u.Phone = phone
//	u.UpdatedAt = time.Now()
//	return nil
//}

//func (u *User) DeleteUserEmail() error {
//	if u.EndedAt != nil {
//		return domain.ErrUserIsDeleted
//	}
//	if u.Email == nil {
//		return domain.ErrInvalidEmail
//	}
//	u.Email = nil
//	u.UpdatedAt = time.Now()
//	return nil
//}

//func (u *User) GetUserContactsByMessengerType(mesTypes []*MessengerType) ([]string, error) {
//	if mesTypes == nil {
//		return nil, errors.New("messenger types are nil")
//	}
//
//	mesContacts := make([]string, 0)
//	for _, mesType := range mesTypes {
//		switch mesType.MesType {
//		case "WhatsApp", "Telegram":
//			mesContacts = append(mesContacts, u.Phone.String())
//		case "Email":
//			if u.Email == nil {
//				return nil, errors.New("user email is nil")
//			}
//			mesContacts = append(mesContacts, string(*u.Email))
//		}
//	}
//	if len(mesContacts) == 0 {
//		return nil, errors.New("user has no contacts")
//	}
//
//	return mesContacts, nil
//}
