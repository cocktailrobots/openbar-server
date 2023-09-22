//go:build !linux

package hardware

import "go.uber.org/zap"

type GpioHardware struct {
}

func (s *GpioHardware) Name() string {
	return "sequent-relay8"
}

func (s *GpioHardware) Close() error {
	return nil
}

func (s *GpioHardware) NumPumps() int {
	return 0
}

func (s *GpioHardware) Pump(idx int, state PumpState) error {
	return nil
}

func (s *GpioHardware) Update(*zap.Logger) {
}

func NewGpioHardware([]int) (*GpioHardware, error) {
	return &GpioHardware{}, nil
}
