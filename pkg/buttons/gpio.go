//go:build !linux

package buttons

import "time"

type GpioButtons struct{}

func NewGpioButtons(pins []int, debounceDur time.Duration, activeHigh, pullUp bool) (*GpioButtons, error) {
	return &GpioButtons{}, nil
}

func (g GpioButtons) NumButtons() int {
	return 0
}

func (g GpioButtons) IsPressed(idx int) bool {
	return false
}

func (g GpioButtons) Update() error {
	return nil
}

func (g GpioButtons) Close() error {
	return nil
}
