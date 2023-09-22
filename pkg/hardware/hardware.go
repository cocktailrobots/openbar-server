package hardware

import (
	"fmt"
	"go.uber.org/zap"
)

type PumpState int

const (
	Off PumpState = iota
	Forward
	Backward
)

func (ps PumpState) String() string {
	switch ps {
	case Off:
		return "Off"
	case Forward:
		return "Forward"
	case Backward:
		return "Backward"
	default:
		return "Unknown"
	}
}

// Hardware interface is the interface for interacting with the pumps and other Barpi hardware
type Hardware interface {
	// Name gets the name of the hardware
	Name() string

	// Close closes the hardware
	Close() error

	// NumPumps gets the number of pumps
	NumPumps() int

	// Pump turns a pump off or on with the given direction
	Pump(idx int, state PumpState) error

	// Update updates the hardware
	Update(logger *zap.Logger)
}

// TurnPumpsOff turns all pumps off
func TurnPumpsOff(h Hardware) error {
	for i := 0; i < h.NumPumps(); i++ {
		if err := h.Pump(i, Off); err != nil {
			return fmt.Errorf("error turning pump %d off: %w", i, err)
		}
	}

	return nil
}
