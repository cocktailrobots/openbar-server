package hardware

import (
	"fmt"
	"log"
	"time"
)

type PumpState int

const (
	Undefined PumpState = iota
	Off
	Forward
	Backward
)

func (ps PumpState) String() string {
	switch ps {
	case Undefined:
		return "Undefined"
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

	// GetReversePin gets the reverse Pin object
	GetReversePin() *ReversePin
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

			if i%3 == 0 {
				hw.update()
				time.Sleep(time.Millisecond)
			}
		}

		hw.update()
	}()

	hw.GetReversePin().SetDirection(Forward)

	start := time.Now()
	onCount := 0
	running := make([]bool, numPumps)
	for i := 0; i < numPumps; i++ {
		if times[i] > 0 {
			if err := hw.pump(i, Forward); err != nil {
				return fmt.Errorf("error turning pump %d on: %w", i, err)
			}

			onCount++
			running[i] = true

			if onCount%3 == 0 {
				hw.update()
				time.Sleep(time.Millisecond)
			}
		}

		if onCount%3 != 0 {
			hw.update()
		}
	}

	for onCount > 0 {
		time.Sleep(5 * time.Millisecond)

		elapsed := time.Since(start)
		onCount = 0
		changes := 0
		for i := 0; i < numPumps; i++ {
			if elapsed <= times[i] {
				onCount++
			} else if running[i] {
				if err := hw.pump(i, Off); err != nil {
					return fmt.Errorf("error turning pump %d off: %w", i, err)
				}

				running[i] = false
				changes++

				if changes%3 == 0 {
					hw.update()
					time.Sleep(time.Millisecond)
				}
			}
		}

		if changes%3 != 0 {
			hw.update()
		}
	}

	return nil
}

type asyncPumpTimes struct {
	times     []time.Time
	direction PumpState
}

func startAsyncHWRoutine(hw Hardware) chan asyncPumpTimes {
	timeCh := make(chan asyncPumpTimes, 32)
	numPumps := hw.NumPumps()

	go func() {
		latestOffTimes := make([]time.Time, numPumps)
		latestDirection := Forward

		currentState := make([]PumpState, numPumps)
		for i := 0; i < numPumps; i++ {
			currentState[i] = Undefined
		}

		for {
			select {
			case apt := <-timeCh:
				if len(apt.times) == numPumps {
					latestOffTimes = apt.times

					if apt.direction != Off {
						latestDirection = apt.direction
					}
				}

			case <-time.After(10 * time.Millisecond):
			}

			now := time.Now()
			for i, t := range latestOffTimes {
				newState := latestDirection
				hw.GetReversePin().SetDirection(newState)

				if t.Before(now) {
					newState = Off
				}

				if newState != currentState[i] {
					if err := hw.Pump(i, newState); err != nil {
						log.Println(err)
					}
					currentState[i] = newState
				}
			}

			//hw.Update()
		}
	}()

	return timeCh
}

type AsyncHWRunner struct {
	hw Hardware
	ch chan asyncPumpTimes
}

func NewAsyncHWRunner(hw Hardware) *AsyncHWRunner {
	return &AsyncHWRunner{
		hw: hw,
		ch: startAsyncHWRoutine(hw),
	}
}

func (ahwr *AsyncHWRunner) RunPumps(direction PumpState, times []time.Duration) error {
	numPumps := ahwr.hw.NumPumps()
	if len(times) != numPumps {
		return fmt.Errorf("expected %d times, but got %d", numPumps, len(times))
	} else if direction == Off {
		return fmt.Errorf("direction cannot be Off")
	} else if direction == Undefined {
		return fmt.Errorf("direction cannot be Undefined")
	}

	apt := asyncPumpTimes{
		times:     make([]time.Time, numPumps),
		direction: direction,
	}

	t := time.Now()
	for i := 0; i < numPumps; i++ {
		apt.times[i] = t.Add(times[i])
	}

	ahwr.ch <- apt
	return nil
}
