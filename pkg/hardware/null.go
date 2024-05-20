package hardware

import (
	"time"
)

type NullHw struct {
	rp *ReversePin
}

func Null(rp *ReversePin) NullHw {
	return NullHw{
		rp: rp,
	}
}

func (nhw NullHw) Name() string {
	return "null"
}

func (nhw NullHw) Close() error {
	return nil
}

func (nhw NullHw) NumPumps() int {
	return 0
}

func (nhw NullHw) Pump(idx int, state PumpState) error {
	return nil
}

func (nhw NullHw) pump(idx int, state PumpState) error {
	return nil
}

func (nhw NullHw) Update() {}
func (nhw NullHw) update() {}

func (nhw NullHw) TimeRun(idx int) time.Duration {
	return 0
}

func (nhw NullHw) RunForTimes(times []time.Duration) error {
	return nil
}

func (nhw NullHw) GetReversePin() *ReversePin {
	return nhw.rp
}
