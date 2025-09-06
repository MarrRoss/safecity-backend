package value_object

type Safety bool

func NewSafety(s bool) (Safety, error) {
	//if err := newSafety.ValidateSafety(); err != nil {
	//	return -1, fmt.Errorf("%w: failed to add safety: %w", err, domain.ErrInvalidZoneSafety)
	//}
	return Safety(s), nil
}

//func (s Safety) ValidateSafety() error {
//	if s < 1 || s > 5 {
//		return errors.New("safety must be no less than 1 and no more than 5")
//	}
//	return nil
//}

func (s Safety) Bool() bool {
	return bool(s)
}
