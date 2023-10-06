//go:build !linux

package hardware

import (
	"time"
)

type SequentRelay8Hardware struct {
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

func (s *SequentRelay8Hardware) RunForTimes(times []time.Duration) error {
	return nil
}

func NewSR8Hardware(expBoardCount int) (*SequentRelay8Hardware, error) {
	return &SequentRelay8Hardware{}, nil
}
