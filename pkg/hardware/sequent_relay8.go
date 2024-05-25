//go:build !linux

package hardware

import (
	"time"
)

type SequentRelay8Hardware struct {
	rp *ReversePin
}

func (s *SequentRelay8Hardware) Name() string {
	return "sequent-relay8"
}

func (s *SequentRelay8Hardware) Close() error {
	return nil
}

func (s *SequentRelay8Hardware) NumPumps() int {
	return 0
}

func (s *SequentRelay8Hardware) Pump(idx int, state PumpState) error {
	return nil
}

func (s *SequentRelay8Hardware) pump(idx int, state PumpState) error {
	return nil
}

func (s *SequentRelay8Hardware) Update() {}
func (s *SequentRelay8Hardware) update() {}

func (s *SequentRelay8Hardware) TimeRun(idx int) time.Duration {
	return 0
}

func (s *SequentRelay8Hardware) RunForTimes(direction PumpState, times []time.Duration) error {
	return nil
}

func (s *SequentRelay8Hardware) GetReversePin() *ReversePin {
	return s.rp
}

func NewSR8Hardware(expBoardCount int, rp *ReversePin) (*SequentRelay8Hardware, error) {
	return &SequentRelay8Hardware{
		rp: rp,
	}, nil
}
