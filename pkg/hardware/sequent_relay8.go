//go:build !linux

package hardware

import "go.uber.org/zap"

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

func (s *SequentRelay8Hardware) Update(*zap.Logger) {
}

func NewSR8Hardware() (*SequentRelay8Hardware, error) {
	return &SequentRelay8Hardware{}, nil
}
