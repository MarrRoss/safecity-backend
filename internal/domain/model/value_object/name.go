package value_object

import (
	"awesomeProjectDDD/internal/domain"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type FullName struct {
	FirstName FirstName
	LastName  LastName
}

func NewFullName(fn FirstName, ln LastName) FullName {
	return FullName{FirstName: fn, LastName: ln}
}

type FirstName string

func NewFirstName(fn string) (FirstName, error) {
	if fn == "" {
		return "", fmt.Errorf("first name is empty: %w", domain.ErrInvalidFirstName)
	}
	var name string
	var err error
	if name, err = ValidateName(fn); err != nil {
		return "", fmt.Errorf("%w: failed to create first name: %w", err, domain.ErrInvalidFirstName)
	}
	return FirstName(name), nil
}

func (fn FirstName) String() string {
	return string(fn)
}

type LastName string

func NewLastName(ln string) (LastName, error) {
	if ln == "" {
		return "", fmt.Errorf("last name is empty: %w", domain.ErrInvalidLastName)
	}
	var name string
	var err error
	if name, err = ValidateName(ln); err != nil {
		return "", fmt.Errorf("%w: failed to create last name: %w", err, domain.ErrInvalidLastName)
	}
	return LastName(name), nil
}

func (ln LastName) String() string {
	return string(ln)
}

type Login string

func NewLogin(login string) (Login, error) {
	if login == "" {
		return "", fmt.Errorf("login is empty: %w", domain.ErrInvalidLogin)
	}
	var name string
	var err error
	if name, err = ValidateName(login); err != nil {
		return "", fmt.Errorf("%w: failed to create login: %w", err, domain.ErrInvalidLogin)
	}
	return Login(name), nil
}

func (login Login) String() string {
	return string(login)
}

//type Patronymic string
//
//func NewPatronymic(pn string) (Patronymic, error) {
//	if pn == "" {
//		return "", fmt.Errorf("patronymic is empty: %w", domain.ErrInvalidPatronymic)
//	}
//	var name string
//	var err error
//	if name, err = ValidateName(pn); err != nil {
//		return "", fmt.Errorf("%w: failed to create patronymic: %w", err, domain.ErrInvalidPatronymic)
//	}
//	return Patronymic(name), nil
//}
//
//func (pn Patronymic) String() string {
//	return string(pn)
//}

type ZoneName string

func NewZoneName(zn string) (ZoneName, error) {
	if zn == "" {
		return "", fmt.Errorf("zone name is empty: %w", domain.ErrInvalidZoneName)
	}
	var name string
	var err error
	if name, err = ValidateName(zn); err != nil {
		return "", fmt.Errorf("%w: failed to create zone name: %w", err, domain.ErrInvalidZoneName)
	}
	return ZoneName(name), nil
}

func (zn ZoneName) String() string {
	return string(zn)
}

type FamilyName string

func NewFamilyName(fn string) (FamilyName, error) {
	if fn == "" {
		return "", fmt.Errorf("family name is empty: %w", domain.ErrInvalidFamilyName)
	}
	var name string
	var err error
	if name, err = ValidateName(fn); err != nil {
		return "", fmt.Errorf("%w: failed to create family name: %w", err, domain.ErrInvalidFamilyName)
	}
	return FamilyName(name), nil
}

func (fn FamilyName) String() string {
	return string(fn)
}

type RoleName string

func NewRoleName(rn string) (RoleName, error) {
	if rn == "" {
		return "", errors.New("role name is empty")
	}
	if err := ValidateRoleName(rn); err != nil {
		return "", fmt.Errorf("failed to create role name: %w", err)
	}
	return RoleName(rn), nil
}

func (rn RoleName) String() string {
	return string(rn)
}

//func ValidateUserName(name string) (string, error) {
//	name = strings.TrimSpace(name)
//
//	if len(name) < 1 || len(name) > 100 {
//		return "", errors.New("name length must be greater than 0 and less than 100")
//	}
//
//	valid := true
//	for _, r := range name {
//		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' {
//			valid = false
//			break
//		}
//	}
//
//	if !valid {
//		return "", errors.New("name must contain only Russian/English letters, spaces and - symbols")
//	}
//
//	runes := []rune(strings.ToLower(name))
//	if len(runes) > 0 {
//		runes[0] = unicode.ToUpper(runes[0])
//	}
//
//	return string(runes), nil
//}

func ValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if len(name) < 1 || len(name) > 100 {
		return "", errors.New("name length must be greater than 0 and less than 100")
	}

	valid := true
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) && r != '-' && r != '_' {
			valid = false
			break
		}
	}

	if !valid {
		return "", errors.New("name must contain only Russian/English letters, numbers, spaces, - and _")
	}

	runes := []rune(strings.ToLower(name))
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes), nil
}

func ValidateRoleName(name string) error {
	if name != "admin" && name != "moderator" && name != "adult" && name != "child" {
		return errors.New("role name must be admin, moderator, adult or child")
	}
	return nil
}
