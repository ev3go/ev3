// Copyright ©2016 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ev3dev provides low level access to the ev3dev control and sensor drivers.
// See documentation at http://www.ev3dev.org/docs/drivers/.
// All functions and methods are safe for concurrent use.
package ev3dev

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	linearPrefix = "linear"
	motorPrefix  = "motor"
	portPrefix   = "port"
	sensorPrefix = "sensor"
)

const (
	// LEDPath is the path to the ev3 LED file system.
	LEDPath = "/sys/class/leds"

	// LegoPortPath is the path to the ev3 lego-port file system.
	LegoPortPath = "/sys/class/lego-port"

	// SensorPath is the path to the ev3 lego-sensor file system.
	SensorPath = "/sys/class/lego-sensor"

	// TachoMotorPath is the path to the ev3 tacho-motor file system.
	TachoMotorPath = "/sys/class/tacho-motor"

	// ServoMotorPath is the path to the ev3 servo-motor file system.
	ServoMotorPath = "/sys/class/servo-motor"

	// DCMotorPath is the path to the ev3 dc-motor file system.
	DCMotorPath = "/sys/class/dc-motor"
)

// These are the subsystem path definitions for all device classes.
const (
	address                   = "address"
	binData                   = "bin_data"
	binDataFormat             = "bin_data_format"
	brightness                = "brightness"
	command                   = "command"
	commands                  = "commands"
	countPerMeter             = "count_per_m"
	countPerRot               = "count_per_rot"
	decimals                  = "decimals"
	direct                    = "direct"
	driverName                = "driver_name"
	dutyCycle                 = "duty_cycle"
	dutyCycleSetpoint         = "duty_cycle_sp"
	fullTravelCount           = "full_travel_count"
	holdPID                   = "hold_pid/"
	holdPIDkd                 = holdPID + kd
	holdPIDki                 = holdPID + ki
	holdPIDkp                 = holdPID + kp
	kd                        = "Kd"
	ki                        = "Ki"
	kp                        = "Kp"
	maxBrightness             = "max_brightness"
	maxPulseSetpoint          = "max_pulse_sp"
	midPulseSetpoint          = "mid_pulse_sp"
	minPulseSetpoint          = "min_pulse_sp"
	maxSpeed                  = "max_speed"
	mode                      = "mode"
	modes                     = "modes"
	numValues                 = "num_values"
	polarity                  = "polarity"
	pollRate                  = "poll_ms"
	position                  = "position"
	positionSetpoint          = "position_sp"
	power                     = "power/"
	powerAutosuspendDelay     = power + "autosuspend_delay_ms"
	powerControl              = power + "control"
	powerRuntimeActiveTime    = power + "runtime_active_time"
	powerRuntimeStatus        = power + "runtime_status"
	powerRuntimeSuspendedTime = power + "runtime_suspended_time"
	rampDownSetpoint          = "ramp_down_sp"
	rampUpSetpoint            = "ramp_up_sp"
	rateSetpoint              = "rate_sp"
	setDevice                 = "set_device"
	speed                     = "speed"
	speedPID                  = "speed_pid/"
	speedPIDkd                = speedPID + kd
	speedPIDki                = speedPID + ki
	speedPIDkp                = speedPID + kp
	speedSetpoint             = "speed_sp"
	state                     = "state"
	status                    = "status"
	stopAction                = "stop_action"
	stopActions               = "stop_actions"
	subsystem                 = "subsystem"
	textValues                = "text_values"
	timeSetpoint              = "time_sp"
	trigger                   = "trigger"
	units                     = "units"
	value                     = "value"
)

// Polarity represent motor polarity states.
type Polarity string

const (
	Normal   Polarity = "normal"
	Inversed Polarity = "inversed"
)

// MotorState is a flag set representing the state of a TachoMotor.
type MotorState uint

const (
	Running MotorState = 1 << iota
	Ramping
	Holding
	Overloaded
	Stalled
)

const (
	running    = "running"
	ramping    = "ramping"
	holding    = "holding"
	overloaded = "overloaded"
	stalled    = "stalled"
)

var motorStateTable = map[string]MotorState{
	running:    Running,
	ramping:    Ramping,
	holding:    Holding,
	overloaded: Overloaded,
	stalled:    Stalled,
}

var motorStates = []string{
	running,
	ramping,
	holding,
	overloaded,
	stalled,
}

// String satisfies the fmt.Stringer interface.
func (f MotorState) String() string {
	const stateMask = Running | Ramping | Holding | Overloaded | Stalled

	var b []byte
	for i, s := range motorStates {
		if f&(1<<uint(i)) != 0 {
			if i != 0 {
				b = append(b, '|')
			}
			b = append(b, s...)
		}
	}

	return string(b)
}

func chomp(b []byte) []byte {
	if b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b
}

// DriverMismatch errors are returned when a device is found that
// does not match the requested driver.
type DriverMismatch struct {
	// Want is the string describing
	// the requested driver.
	Want string

	// Have is the string describing
	// the driver present on the device.
	Have string
}

func (e DriverMismatch) Error() string {
	return fmt.Sprintf("ev3dev: mismatched driver names: want %q but have %q", e.Want, e.Have)
}

// Device is an ev3dev API device.
type Device interface {
	// Path returns the sysfs path
	// for the device type.
	Path() string

	// Type returns the type of the
	// device, one of "linear", "motor",
	// "port" or "sensor".
	Type() string

	fmt.Stringer
}

// AddressOf returns the port address of the Device.
func AddressOf(d Device) (string, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf(d.Path()+"/%s/"+address, d))
	if err != nil {
		return "", fmt.Errorf("ev3dev: failed to read %s address: %v", d.Type(), err)
	}
	return string(chomp(b)), err
}

// DriverFor returns the driver name for the Device.
func DriverFor(d Device) (string, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf(d.Path()+"/%s/"+driverName, d))
	if err != nil {
		return "", fmt.Errorf("ev3dev: failed to read %s driver name: %v", d.Type(), err)
	}
	return string(chomp(b)), err
}

// deviceIDFor returns the id for the given ev3 port name and driver of the Device.
// If the driver does not match the driver string, an id for the device is returned
// with a DriverMismatch error.
// If port is empty, the first device satisfying the driver name is returned.
func deviceIDFor(port, driver string, d Device) (int, error) {
	f, err := os.Open(d.Path())
	if err != nil {
		return -1, err
	}
	devices, err := f.Readdirnames(0)
	f.Close()
	if err != nil {
		return -1, fmt.Errorf("ev3dev: could not get devices for %s: %v", d.Path(), err)
	}

	portBytes := []byte(port)
	driverBytes := []byte(driver)
	for _, device := range devices {
		if !strings.HasPrefix(device, d.Type()) {
			continue
		}
		id, err := strconv.Atoi(strings.TrimPrefix(device, d.Type()))
		if err != nil {
			return -1, fmt.Errorf("ev3dev: could not parse id from device name %q: %v", device, err)
		}

		if port == "" {
			path := filepath.Join(d.Path(), device, driverName)
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return -1, fmt.Errorf("ev3dev: could not read driver name %s: %v", path, err)
			}
			if !bytes.Equal(driverBytes, chomp(b)) {
				continue
			}
			return id, nil
		}

		path := filepath.Join(d.Path(), device, address)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return -1, fmt.Errorf("ev3dev: could not read address %s: %v", path, err)
		}
		if !bytes.Equal(portBytes, chomp(b)) {
			continue
		}
		path = filepath.Join(d.Path(), device, driverName)
		b, err = ioutil.ReadFile(path)
		if err != nil {
			return -1, fmt.Errorf("ev3dev: could not read driver name %s: %v", path, err)
		}
		if !bytes.Equal(driverBytes, chomp(b)) {
			err = DriverMismatch{Want: driver, Have: string(b)}
		}
		return id, err
	}

	if port != "" {
		return -1, fmt.Errorf("ev3dev: could not find device for driver %q on port %s", driver, port)
	}
	return -1, fmt.Errorf("ev3dev: could not find device for driver %q", driver)
}
