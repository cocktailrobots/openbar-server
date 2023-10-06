package buttons

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
)

type GpioButtons struct {
	lines []*gpiod.Line
}

func NewGpioButtons(pins []int, debounceDur time.Duration, activeLow, pullUp bool) (*GpioButtons, error) {
	options := []gpiod.LineReqOption{gpiod.AsInput}
	if activeLow {
		options = append(options, gpiod.AsActiveLow)
	} else {
		options = append(options, gpiod.AsActiveHigh)
	}

	if pullUp {
		options = append(options, gpiod.WithPullUp)
	} else {
		options = append(options, gpiod.WithPullDown)
	}

	if debounceDur > 0 {
		options = append(options, gpiod.WithDebounce(debounceDur))
	}

	var lines []*gpiod.Line
	for _, pin := range pins {
		l, err := gpiod.RequestLine("gpiochip0", int(pin), options...)
		if err != nil {
			return nil, fmt.Errorf("error requesting line %d as input: %w", pin, err)
		}

		lines = append(lines, l)
	}

	return &GpioButtons{
		lines: lines,
	}, nil
}

func (g GpioButtons) NumButtons() int {
	return len(g.lines)
}

func (g GpioButtons) IsPressed(idx int) bool {
	val, err := g.lines[idx].Value()
	if err != nil {
		panic(err)
	}

	return val == 1
}

func (g GpioButtons) Update() error {
	return nil
}

func (g GpioButtons) Close() error {
	for _, l := range g.lines {
		l.Close()
	}

	return nil
}
