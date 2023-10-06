package hardware

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/util/ncurses"
	"github.com/gbin/goncurses"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type stateChange struct {
	state     PumpState
	changedAt time.Time
}

// DebugHardware is the hardware implementation for debugging
type DebugHardware struct {
	mu          *sync.Mutex
	numPumps    int
	outFilePath string

	initialOut *os.File
	initialErr *os.File

	win   *goncurses.Window
	table *ncurses.Table

	state    []stateChange
	runTimes []time.Duration
}

// NewDebugHardware creates a new DebugHardware
func NewDebugHardware(numPumps int, outFilePath string) (*DebugHardware, error) {
	f, err := os.Create(outFilePath)
	if err != nil {
		return nil, fmt.Errorf("error creating file %s: %v", outFilePath, err)
	}

	initialOut := os.Stdout
	initialErr := os.Stderr
	os.Stderr = f
	os.Stdout = f

	win, err := goncurses.Init()
	if err != nil {
		os.Stdout = initialOut
		os.Stderr = initialErr
		f.Close()

		return nil, fmt.Errorf("error initializing ncurses: %w", err)
	}

	now := time.Now()
	initialState := make([]stateChange, numPumps)
	for i := 0; i < numPumps; i++ {
		initialState[i] = stateChange{
			state:     Off,
			changedAt: now,
		}
	}

	return &DebugHardware{
		mu:          &sync.Mutex{},
		numPumps:    numPumps,
		outFilePath: outFilePath,
		initialOut:  initialOut,
		initialErr:  initialErr,
		state:       initialState,
		win:         win,
		table:       ncurses.NewTable([]int{2, 8, 10}),
		runTimes:    make([]time.Duration, numPumps),
	}, nil
}

// Name gets the name of the hardware
func (h *DebugHardware) Name() string {
	return "Debug"
}

// Close closes the hardware
func (h *DebugHardware) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	goncurses.End()

	f := os.Stdout
	os.Stdout = h.initialOut
	os.Stderr = h.initialErr

	err := f.Close()
	if err != nil {
		return fmt.Errorf("error closing file %s: %w", h.outFilePath, err)
	}

	log.Printf("Debug Hardware closed. Output written to %s", h.outFilePath)
	return nil
}

// NumPumps gets the number of pumps
func (h *DebugHardware) NumPumps() int {
	return h.numPumps
}

// Pump turns a pump off or on with the given direction
func (h *DebugHardware) Pump(idx int, state PumpState) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.pump(idx, state)
}

func (h *DebugHardware) pump(idx int, state PumpState) error {
	if idx < 0 || idx >= h.numPumps {
		return fmt.Errorf("invalid pump index %d", idx)
	}

	currState := h.state[idx].state
	newState := state
	if currState != newState {
		now := time.Now()
		if currState == Forward {
			h.runTimes[idx] += now.Sub(h.state[idx].changedAt)
		}

		h.state[idx].state = state
		h.state[idx].changedAt = now
	}
	return nil
}

// Update updates the hardware
func (h *DebugHardware) Update() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.render()
}

// update updates the hardware without locking for internal use
func (h *DebugHardware) update() {
	h.render()
}

func (h *DebugHardware) render() {
	h.win.Clear()
	h.win.MovePrint(0, 0, "Debug Hardware")
	h.win.MovePrint(1, 0, "Pump State:")

	now := time.Now()
	rows := [][]string{
		{"#", "State", "Elapsed"},
	}

	for i := range h.state {
		s := h.state[i]
		elapsed := now.Sub(s.changedAt).Seconds()
		rows = append(rows, []string{
			strconv.FormatInt(int64(i), 10),
			s.state.String(),
			fmt.Sprintf("%0.03f", elapsed),
		})
	}

	h.table.Render(h.win, 2, 2, rows)
	h.win.Refresh()
}

// TimeRun returns the total time the pump has been run for since the program started
func (h *DebugHardware) TimeRun(idx int) time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()

	if idx < 0 || idx >= h.numPumps {
		panic(fmt.Errorf("invalid pump index %d", idx))
	}

	return h.runTimes[idx]
}

// RunForTimes runs the pumps for the given times
func (h *DebugHardware) RunForTimes(times []time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return runForTimes(h, times)
}
