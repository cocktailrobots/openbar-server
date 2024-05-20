package hardware

import (
	"sync"

	"github.com/warthog618/gpiod"

	cfg "github.com/cocktailrobots/openbar-server/pkg/config"
)

type ReversePin struct {
	mu         *sync.Mutex
	pin        int
	forwardVal int
	backVal    int
	currentVal int
	line       *gpiod.Line
}

func NewReversePin(config *cfg.ReversePinConfig) (*ReversePin, error) {
	if cfg == nil || config.Pin == -1 {
		return &ReversePin{
			pin:        -1,
			forwardVal: 0,
			backVal:    1,
		}, nil
	}

	l, err := gpiod.RequestLine("gpiochip0", pin, gpiod.AsOutput(0))
	if err != nil {
		return nil, fmt.Errorf("error requesting line %d as output: %w", pin, err)
	}

	var backVal int
	if forwardVal == 0 {
		backVal = 1
	}

	return &ReversePin{
		pin:        pin,
		line:       l,
		forwardVal: forwardVal,
		backVal:    backVal,
	}
}

func (rp *ReversePin) SetDirection(direction PumpState) error {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	var val int
	if direction == Forward {
		val = rp.forwardVal
	} else {
		val = rp.backVal
	}

	if rp.line != nil && val != rp.currentVal {
		err := rp.line.SetValue(val)
		if err != nil {
			return fmt.Errorf("error setting value %d on line %d: %w", val, rp.pin, err)
		}
	}

	rp.currentVal = val
	return nil
}
