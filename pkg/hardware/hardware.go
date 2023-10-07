package hardware

import (
	"fmt"
	"log"
	"time"
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

	// pump is a lock free version of Pump for package internal use
	pump(idx int, state PumpState) error

	// Update updates the hardware
	Update()

	// update updates the hardware without locking for internal use
	update()

	// TimeRun returns the total time the pump has been run for since the program started
	TimeRun(idx int) time.Duration

	// RunForTimes runs the pumps for the given times
	RunForTimes(times []time.Duration) error
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

func runForTimes(hw Hardware, times []time.Duration) error {
	numPumps := hw.NumPumps()
	if len(times) != numPumps {
		return fmt.Errorf("expected %d times, but got %d", numPumps, len(times))
	}

	defer func() {
		for i := 0; i < numPumps; i++ {
			err := hw.pump(i, Off)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	start := time.Now()
	onCount := 0
	for i := 0; i < numPumps; i++ {
		if times[i] > 0 {
			if err := hw.pump(i, Forward); err != nil {
				return fmt.Errorf("error turning pump %d on: %w", i, err)
			}

			onCount++
		}
	}
	hw.update()

	for onCount > 0 {
		elapsed := time.Since(start)
		onCount = 0
		for i := 0; i < numPumps; i++ {
			if elapsed <= times[i] {
				onCount++
			} else {
				if err := hw.pump(i, Off); err != nil {
					return fmt.Errorf("error turning pump %d to state %s: %w", i, Off.String(), err)
				}
			}
		}

		hw.update()

		if onCount > 0 {
			time.Sleep(5 * time.Millisecond)
		}
	}

	return nil
}
