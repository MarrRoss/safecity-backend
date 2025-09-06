package value_object

//type Password string
//
//func NewPassword(p string) (Password, error) {
//	if p == "" {
//		return "", fmt.Errorf("password is empty: %w", domain.ErrInvalidPassword)
//	}
//	newPassword := Password(p)
//	if err := newPassword.ValidatePassword(); err != nil {
//		return "", fmt.Errorf("%w: %w", err, domain.ErrInvalidPassword)
//	}
//	return newPassword, nil
//}
//
//func (p Password) ValidatePassword() error {
//	if len(p) < 8 {
//		return errors.New("password must be at least 8 characters long")
//	}
//
//	hasUpperCase := false
//	for _, char := range p {
//		if unicode.IsUpper(char) {
//			hasUpperCase = true
//			break
//		}
//	}
//	if !hasUpperCase {
//		return errors.New("password must contain at least one uppercase letter")
//	}
//
//	hasDigit := false
//	for _, char := range p {
//		if unicode.IsDigit(char) {
//			hasDigit = true
//			break
//		}
//	}
//	if !hasDigit {
//		return errors.New("password must contain at least one digit")
//	}
//	return nil
//}
//
//func (p Password) String() string {
//	return string(p)
//}
