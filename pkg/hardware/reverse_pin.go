//go:build !linux

package hardware

import (
	cfg "github.com/cocktailrobots/openbar-server/pkg/config"
	"sync"
)

type ReversePin struct {
	mu         *sync.Mutex
	forwardVal int
	currentVal int
}

func NewReversePin(config *cfg.ReversePinConfig) (*ReversePin, error) {
	fv := 0
	if config != nil && config.ForwardHigh {
		fv = 1
	}

	return &ReversePin{
		mu:         &sync.Mutex{},
		forwardVal: fv,
	}, nil
}

func (rp *ReversePin) SetDirection(direction PumpState) error {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	rp.currentVal = rp.forwardVal
	if direction == Backward {
		rp.currentVal = 1 - rp.forwardVal
	}

	return nil
}

func (rp *ReversePin) Value() int {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	return rp.currentVal
}
