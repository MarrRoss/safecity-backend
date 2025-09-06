package value_object

import (
	"awesomeProjectDDD/internal/domain"
	"errors"
	"fmt"
	"regexp"
)

type Email string

func NewEmail(e string) (Email, error) {
	newEmail := Email(e)
	if err := newEmail.ValidateEmail(); err != nil {
		return "", fmt.Errorf("%w: %w", err, domain.ErrInvalidEmail)
	}
	return newEmail, nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) ValidateEmail() error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(emailRegex, string(e))
	if !match {
		return errors.New("invalid format")
	}

	return nil
}
