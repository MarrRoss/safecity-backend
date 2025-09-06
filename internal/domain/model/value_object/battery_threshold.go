package value_object

import (
	"awesomeProjectDDD/internal/domain"
)

type BatteryThreshold int

func NewBatteryThreshold(battery int) (BatteryThreshold, error) {
	newBattery := BatteryThreshold(battery)
	if err := newBattery.ValidateBatteryThreshold(); err != nil {
		return -1, err
	}
	return newBattery, nil
}

func (b BatteryThreshold) ValidateBatteryThreshold() error {
	if b < 0 || b > 100 {
		return domain.ErrInvalidBattery
		//return errors.New("battery threshold must be in the range [0;100]")
	}
	return nil
}

func (b BatteryThreshold) Int() int {
	return int(b)
}
