package hardware

import (
	"sync"
	"time"
)

// Hardware interface is the interface for interacting with the pumps and other Barpi hardware
type TestHardware struct {
	mu        *sync.Mutex
	numPumps  int
	state     []PumpState
	runTimes  []time.Duration
	changedAt []time.Time
}

func NewTestHardware(numPumps int) *TestHardware {
	return &TestHardware{
		mu:        &sync.Mutex{},
		numPumps:  numPumps,
		state:     make([]PumpState, 8),
		runTimes:  make([]time.Duration, 8),
		changedAt: make([]time.Time, 8),
	}
}

func (thw *TestHardware) ResetRuntimes() {
	thw.runTimes = make([]time.Duration, thw.numPumps)
}

func (thw *TestHardware) SetRuntime(idx int, runtime time.Duration) {
	thw.mu.Lock()
	defer thw.mu.Unlock()

	thw.runTimes[idx] = runtime
}

func (thw *TestHardware) Name() string {
	return "test"
}

func (thw *TestHardware) Close() error {
	return nil
}

func (thw *TestHardware) NumPumps() int {
	return thw.numPumps
}

func (thw *TestHardware) Pump(idx int, state PumpState) error {
	thw.mu.Lock()
	defer thw.mu.Unlock()

	return thw.pump(idx, state)
}

func (thw *TestHardware) pump(idx int, state PumpState) error {
	currState := thw.state[idx]
	if currState == state {
		return nil
	}

	if currState == Forward {
		thw.runTimes[idx] += time.Now().Sub(thw.changedAt[idx])
	}

	thw.state[idx] = state
	thw.changedAt[idx] = time.Now()

	return nil
}

func (thw *TestHardware) Update() {
	thw.mu.Lock()
	defer thw.mu.Unlock()

	thw.update()
}

func (thw *TestHardware) update() {}

func (thw *TestHardware) TimeRun(idx int) time.Duration {
	thw.mu.Lock()
	defer thw.mu.Unlock()

	return thw.runTimes[idx]
}

func (thw *TestHardware) RunForTimes(times []time.Duration) error {
	thw.mu.Lock()
	defer thw.mu.Unlock()

	return runForTimes(thw, times)
}
