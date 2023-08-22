package hardware

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/util/ncurses"
	"github.com/gbin/goncurses"
	"go.uber.org/zap"
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

type stateChanges struct {
	mu      *sync.Mutex
	changes []stateChange
}

func NewStateChanges(initialVals []stateChange) stateChanges {
	return stateChanges{
		mu:      &sync.Mutex{},
		changes: initialVals,
	}
}

func (s *stateChanges) GetState(pumpNum int) stateChange {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.changes[pumpNum]
}

func (s *stateChanges) SetState(pumpNum int, state PumpState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.changes[pumpNum] = stateChange{
		state:     state,
		changedAt: time.Now(),
	}
}

// DebugHardware is the hardware implementation for debugging
type DebugHardware struct {
	numPumps    int
	outFilePath string

	initialOut *os.File
	initialErr *os.File

	win   *goncurses.Window
	table *ncurses.Table

	state stateChanges
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
		numPumps:    numPumps,
		outFilePath: outFilePath,
		initialOut:  initialOut,
		initialErr:  initialErr,
		state:       NewStateChanges(initialState),
		win:         win,
		table:       ncurses.NewTable([]int{2, 8, 10}),
	}, nil
}

// Name gets the name of the hardware
func (h *DebugHardware) Name() string {
	return "Debug"
}

// Close closes the hardware
func (h *DebugHardware) Close() error {
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
	if idx < 0 || idx >= h.numPumps {
		return fmt.Errorf("invalid pump index %d", idx)
	}

	if h.state.GetState(idx).state != state {
		h.state.SetState(idx, state)
	}
	return nil
}

// Update updates the hardware
func (h *DebugHardware) Update(*zap.Logger) {
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

	for i := range h.state.changes {
		s := h.state.GetState(i)
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
