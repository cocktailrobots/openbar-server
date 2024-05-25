package hardware

import (
	"fmt"
	"sync"
	"time"

	"github.com/warthog618/gpiod"
)

var _ Hardware = &GpioHardware{}

type pump struct {
	pin       int
	line      *gpiod.Line
	state     PumpState
	updatedAt time.Time
}

func (p pump) String() string {
	return fmt.Sprintf("%s for %f seconds", p.state.String(), time.Since(p.updatedAt).Seconds())
}

// GpioHardware is the hardware implementation for the Raspberry Pi GPIO pins
type GpioHardware struct {
	mu       *sync.Mutex
	pumps    []pump
	runTimes []time.Duration
	rp       *ReversePin
}

func NewGpioHardware(pins []int, rp *ReversePin) (*GpioHardware, error) {
	var pumps []pump
	for _, pin := range pins {
		l, err := gpiod.RequestLine("gpiochip0", pin, gpiod.AsOutput(0))
		if err != nil {
			return nil, fmt.Errorf("error requesting line %d as output: %w", pin, err)
		}

		pumps = append(pumps, pump{
			pin:       pin,
			line:      l,
			state:     Off,
			updatedAt: time.Now(),
		})
	}

	return &GpioHardware{
		mu:       &sync.Mutex{},
		pumps:    pumps,
		runTimes: make([]time.Duration, len(pumps)),
		rp:       rp,
	}, nil
}

func (g *GpioHardware) Name() string {
	return "GPIO"
}

func (g *GpioHardware) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	var firstErr error
	for _, p := range g.pumps {
		err := p.line.Close()
		if firstErr == nil && err != nil {
			firstErr = fmt.Errorf("error closing line %d: %w", p.pin, err)
		}
	}

	return firstErr
}

func (g *GpioHardware) NumPumps() int {
	return len(g.pumps)
}

func (g *GpioHardware) Pump(idx int, state PumpState) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.pump(idx, state)
}

func (g *GpioHardware) pump(idx int, state PumpState) error {
	if idx < 0 || idx >= g.NumPumps() {
		return fmt.Errorf("invalid pump index %d", idx)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()
	if g.pumps[idx].state == Forward && state != Forward {
		g.runTimes[idx] += now.Sub(g.pumps[idx].updatedAt)
	}

	p := g.pumps[idx]
	var err error
	switch state {
	case Off:
		err = p.line.SetValue(0)
	case Forward:
		err = p.line.SetValue(1)
	case Backward:
		err = fmt.Errorf("Not implemented")
	default:
		err = fmt.Errorf("unknown state %d", state)
	}

	if err != nil {
		return fmt.Errorf("error setting line %d to state %s: %w", p.pin, state, err)
	}

	p.state = state
	p.updatedAt = now
	return nil
}

func (g *GpioHardware) Update() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.update()
}

func (g *GpioHardware) update() {
}

func (g *GpioHardware) TimeRun(idx int) time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()

	if idx < 0 || idx >= g.NumPumps() {
		panic(fmt.Errorf("invalid pump index %d", idx))
	}

	return g.runTimes[idx]
}

func (g *GpioHardware) RunForTimes(direction PumpState, times []time.Duration) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return runForTimes(g, direction, times)
}

func (g *GpioHardware) GetReversePin() *ReversePin {
	return g.rp
}
