package value_object

//type BirthDate time.Time
//
//func NewBirthDate(bd time.Time) (BirthDate, error) {
//	if bd.IsZero() {
//		return BirthDate{}, fmt.Errorf("%w: birth date is empty", domain.ErrInvalidBirthDate)
//	}
//	newBirthDate := BirthDate(bd)
//	if err := newBirthDate.ValidateAge(); err != nil {
//		return BirthDate{}, fmt.Errorf("%w: invalid format: %w", domain.ErrInvalidBirthDate, err)
//	}
//	return newBirthDate, nil
//}
//
//func NewBirthDateFromString(bd string) (BirthDate, error) {
//	if bd == "" {
//		return BirthDate{}, fmt.Errorf("%w: birth date is empty", domain.ErrInvalidBirthDate)
//	}
//	layout := "2006-02-01"
//	parsedDate, err := time.Parse(layout, bd)
//	if err != nil {
//		return BirthDate{}, fmt.Errorf("%w: invalid format, expected DD-MM-YYYY", domain.ErrInvalidBirthDate)
//	}
//	return NewBirthDate(parsedDate)
//}
//
//func NewBirthDateFromTime(bd time.Time) (BirthDate, error) {
//	if bd.IsZero() {
//		return BirthDate{}, fmt.Errorf("birth date is empty: %w", domain.ErrInvalidBirthDate)
//	}
//	formatted := bd.Format("2006-01-02")
//	parsedDate, err := time.Parse("2006-01-02", formatted)
//	if err != nil {
//		return BirthDate{}, fmt.Errorf("failed to parse formatted date: %w", domain.ErrInvalidBirthDate)
//	}
//	return NewBirthDate(parsedDate)
//}
//
//func (bd BirthDate) String() string {
//	return time.Time(bd).Format("2006-02-01")
//}
//
//func (bd BirthDate) Time() time.Time {
//	return time.Time(bd)
//}
//
//func (bd BirthDate) ValidateAge() error {
//	age := time.Since(time.Time(bd)).Hours() / 24 / 365
//	if age < 5 || age > 130 {
//		return fmt.Errorf("age must be between 5 and 130 years")
//	}
//	return nil
//}
