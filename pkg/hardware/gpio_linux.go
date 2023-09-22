package hardware

import (
	"fmt"
	"strings"
	"time"

	"github.com/warthog618/gpiod"
	"go.uber.org/zap"
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
	pumps []pump
}

func NewGpioHardware(pins []int) (*GpioHardware, error) {
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
		pumps: pumps,
	}, nil
}

func (g GpioHardware) Name() string {
	return "GPIO"
}

func (g GpioHardware) Close() error {
	var firstErr error
	for _, p := range g.pumps {
		err := p.line.Close()
		if firstErr == nil && err != nil {
			firstErr = fmt.Errorf("error closing line %d: %w", p.pin, err)
		}
	}

	return firstErr
}

func (g GpioHardware) NumPumps() int {
	return len(g.pumps)
}

func (g GpioHardware) Pump(idx int, state PumpState) error {
	if idx < 0 || idx >= g.NumPumps() {
		return fmt.Errorf("invalid pump index %d", idx)
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
	p.updatedAt = time.Now()
	return nil
}

func (g GpioHardware) Update(logger *zap.Logger) {
	var strs []string
	for i, p := range g.pumps {
		strs = append(strs, fmt.Sprintf(`"pump_%d": "%s"`, i, p.String()))
	}

	logger.Info("GpioHardware Update", zap.String("state", strings.Join(strs, ", ")))
}