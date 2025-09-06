package value_object

//type Phone string
//
//func NewPhone(p string) (Phone, error) {
//	if p == "" {
//		return "", fmt.Errorf("phone is empty: %w", domain.ErrInvalidPhone)
//	}
//	phone, err := ValidatePhone(p)
//	if err != nil {
//		return "", fmt.Errorf("%w: %w", err, domain.ErrInvalidPhone)
//	}
//	return Phone(phone), nil
//}
//
//func (p Phone) isAllDigits(s string) bool {
//	if len(s) == 0 {
//		return false
//	}
//	for _, ch := range s {
//		if !unicode.IsDigit(ch) {
//			return false
//		}
//	}
//	return true
//}
//
//func (p Phone) String() string {
//	return string(p)
//}
//
//func ValidatePhone(p string) (string, error) {
//	trimmed := strings.TrimSpace(p)
//
//	if len(trimmed) > 0 && trimmed[0] == '+' {
//		trimmed = trimmed[1:]
//	}
//
//	for _, ch := range trimmed {
//		if ch < '0' || ch > '9' {
//			return "", errors.New("phone must contain only digits")
//		}
//	}
//
//	if len(trimmed) != 11 {
//		return "", errors.New("phone must be 11 characters long")
//	}
//
//	if trimmed[0] != '7' && trimmed[0] != '8' {
//		return "", errors.New("phone must start with 7 or 8")
//	}
//
//	if trimmed[0] == '8' {
//		trimmed = "7" + trimmed[1:]
//	}
//
//	return trimmed, nil
//}
