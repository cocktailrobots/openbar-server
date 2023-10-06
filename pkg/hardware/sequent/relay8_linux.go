package sequent

import (
	"fmt"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"log"
	"strings"
)

type RelayState bool
type Relay8States [8]RelayState

func (r Relay8States) toByte() byte {
	var state byte
	for i := range r {
		if r[i] {
			state |= relayMaskRemap[i]
		}
	}

	return state
}

func fromByte(state byte) Relay8States {
	var r Relay8States
	for i := range r {
		r[i] = state&relayMaskRemap[i] != 0
	}

	return r
}

func (r Relay8States) Set(idx int, state RelayState) Relay8States {
	r[idx] = state
	return r
}

func (r Relay8States) Get(idx int) RelayState {
	return r[idx]
}

func (r Relay8States) Equal(other Relay8States) bool {
	for i := range r {
		if r[i] != other[i] {
			return false
		}
	}

	return true
}

const (
	relayOn  = true
	relayOff = false

	relay8HwI2cBaseAddr    = 0x38
	relay8HwI2cAltBaseAddr = 0x20
	relay8CfgRegAddr       = 0x03
	relay8OutportRegAddr   = 0x01
	relay8InportRegAddr    = 0x00
)

var relayMaskRemap = []byte{0x01, 0x04, 0x40, 0x10, 0x20, 0x80, 0x08, 0x02}

var relayChRemap = []int{0, 2, 6, 4, 5, 7, 3, 1}

// verbose controls whether or not to print verbose logging
var verbose = false

func init() {
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
}

func SetVerboseLogging(enabled bool) {
	verbose = enabled

	if enabled {
		logger.ChangePackageLogLevel("i2c", logger.DebugLevel)
	} else {
		logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	}
}

func logln(a ...any) {
	if verbose {
		log.Println(a...)
	}
}

func logf(f string, a ...any) {
	if verbose {
		log.Printf(f, a...)
	}
}

func bytesAsHex(buff []byte) string {
	var sb strings.Builder
	for i := range buff {
		sb.WriteString(fmt.Sprintf("%02X ", buff[i]))
	}
	return sb.String()
}

func boardCheck(hwAdd byte) (bool, error) {
	hwAdd ^= 0x07
	dev, err := i2c.NewI2C(hwAdd, 1)
	if err != nil {
		return false, fmt.Errorf("error creating I2C device: %w", err)
	}
	defer DeinitBoard(dev)

	buff, n, err := dev.ReadRegBytes(relay8CfgRegAddr, 8)
	if err != nil {
		return false, fmt.Errorf("error reading from device: %w", err)
	}

	logf("board check read %d bytes: %s", n, bytesAsHex(buff[:n]))
	return true, nil
}

func DeinitBoard(dev *i2c.I2C) {
	err := dev.Close()
	if err != nil && verbose {
		log.Println("error closing device:", err)
	}
}

func InitBoard(stack byte) (*i2c.I2C, error) {
	if stack < 0 || stack > 7 {
		panic(fmt.Errorf("stack %d is not between 0 and 7", stack))
	}

	addr := (stack + relay8HwI2cBaseAddr) ^ 0x07
	dev, err := i2c.NewI2C(addr, 1)
	if err != nil {
		return nil, fmt.Errorf("error creating I2C device with addr 0x%x: %w", addr, err)
	}

	buff, n, err := dev.ReadRegBytes(relay8CfgRegAddr, 1)
	if err != nil {
		dev.Close()

		addr = (stack + relay8HwI2cAltBaseAddr) ^ 0x07
		dev, err = i2c.NewI2C(addr, 1)
		if err != nil {
			return nil, fmt.Errorf("error creating I2C device with addr 0x%x: %w", addr, err)
		}

		buff, n, err = dev.ReadRegBytes(relay8CfgRegAddr, 1)
		if err != nil {
			dev.Close()
			return nil, fmt.Errorf("error reading from device: %w", err)
		}
	}

	if n != 1 {
		dev.Close()
		return nil, fmt.Errorf("error, only %d of %d bytes were read", n, 1)
	}

	if buff[0] != 0 { //non initialized I/O Expander
		// make all I/O pins output
		if err := dev.WriteRegU8(relay8CfgRegAddr, 0); err != nil {
			dev.Close()
			return nil, fmt.Errorf("error writing to device to setup io pins: %w", err)
		}

		// put all pins in 0-logic state
		if err := dev.WriteRegU8(relay8OutportRegAddr, 0); err != nil {
			dev.Close()
			return nil, fmt.Errorf("error writing to device to initialize io pins to 0: %w", err)
		}
	}

	return dev, nil
}

func writeRelay8(dev *i2c.I2C, states byte) error {
	logf("writing states 0x%02x", states)
	if err := dev.WriteRegU8(relay8OutportRegAddr, states); err != nil {
		return fmt.Errorf("error writing to device: %w", err)
	}

	return nil
}

func readRelay8(dev *i2c.I2C) (byte, error) {
	buff, _, err := dev.ReadRegBytes(relay8InportRegAddr, 1)
	if err != nil {
		return 0, fmt.Errorf("error reading from device: %w", err)
	}

	return buff[0], nil
}

func UpdateBoard(dev *i2c.I2C, states Relay8States, attempts int) error {
	desired := states.toByte()

	for i := attempts - 1; i >= 0; i++ {
		err := writeRelay8(dev, desired)
		if err != nil {
			logf("error writing to device: %s", err.Error())
			if i == 0 {
				return fmt.Errorf("error writing to device: %w", err)
			}
		}

		read, err := readRelay8(dev)
		if err != nil {
			logf("error reading from device: %s", err.Error())
			if i == 0 {
				return fmt.Errorf("error reading from device: %w", err)
			}
		}

		if read == desired {
			return nil
		}
	}

	return fmt.Errorf("never achieved desired state 0x%02x", desired)
}
