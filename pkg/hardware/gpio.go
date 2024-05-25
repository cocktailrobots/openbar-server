//go:build !linux

package hardware

import (
	"time"
)

type GpioHardware struct {
	rp *ReversePin
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

func (s *GpioHardware) pump(idx int, state PumpState) error {
	return nil
}

func (s *GpioHardware) Update() {
	s.update()
}

func (s *GpioHardware) update() {
}

func (s *GpioHardware) TimeRun(idx int) time.Duration {
	return 0
}

func (s *GpioHardware) RunForTimes(direction PumpState, times []time.Duration) error {
	return nil
}

func (s *GpioHardware) GetReversePin() *ReversePin {
	return s.rp
}

func NewGpioHardware(pins []int, rp *ReversePin) (*GpioHardware, error) {
	return &GpioHardware{
		rp: rp,
	}, nil
}
