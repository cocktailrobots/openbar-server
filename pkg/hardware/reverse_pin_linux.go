package hardware

import (
	"fmt"
	"log"
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
	if config == nil || config.Pin == -1 {
		return &ReversePin{
			mu:         &sync.Mutex{},
			pin:        -1,
			forwardVal: 0,
			backVal:    1,
		}, nil
	}

	l, err := gpiod.RequestLine("gpiochip0", config.Pin, gpiod.AsOutput(1))
	if err != nil {
		return nil, fmt.Errorf("error requesting line %d as output: %w", config.Pin, err)
	}

	var forwardVal int
	var backVal int = 1
	if config.ForwardHigh {
		backVal = 0
		forwardVal = 1
	}

	return &ReversePin{
		mu:         &sync.Mutex{},
		pin:        config.Pin,
		line:       l,
		forwardVal: forwardVal,
		backVal:    backVal,
	}, nil
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
		log.Printf("setting pin %d to %d", rp.pin, val)
		err := rp.line.SetValue(val)
		if err != nil {
			return fmt.Errorf("error setting value %d on line %d: %w", val, rp.pin, err)
		}
	}

	rp.currentVal = val
	return nil
}

func (rp *ReversePin) Value() int {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	return rp.currentVal
}
