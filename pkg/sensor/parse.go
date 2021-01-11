package sensor

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// data sheet for open close sensor
// https://wings.twelite.info/how-to-use/parent-mode/receive-message/app_pal#sensparu

const (
	ModeCloser = 0x81 + iota
	ModeEnv
	ModeMotion
	ModeNotification
)

const (
	NoMagnet = iota
	CloseNPole
	CloseSPole
	PeriodicTransmission = 0x80
)

type Sensor interface {
	parse(header commonHeader, data []byte) (Sensor, error)
}

type commonHeader struct {
	SerialID     string
	DeviceID     string
	LQI          byte
	PowerVoltage int
}

var _ Sensor = CloserSensor{}

type CloserSensor struct {
	commonHeader
	State byte
}

func (s CloserSensor) parse(header commonHeader, data []byte) (Sensor, error) {
	s.commonHeader = header
	s.State = data[31]
	return s, nil
}

func Parse(data string) (Sensor, error) {
	if err := validateData(data); err != nil {
		return nil, err
	}
	h, err := hex.DecodeString(data[1:])
	if err != nil {
		return nil, err
	}
	if err := verifyChecksum(h); err != nil {
		return nil, err
	}
	header := commonHeader{
		SerialID:     data[15:23],
		DeviceID:     data[23:25],
		LQI:          h[4],
		PowerVoltage: int(h[19])<<8 + int(h[20]),
	}
	var sensor Sensor
	switch mode := h[13]; mode {
	case ModeCloser:
		sensor = CloserSensor{}
	case ModeEnv:
		return nil, errors.New("environment sensor is not implemented")
	case ModeMotion:
		return nil, errors.New("motion sensor is not implemented")
	case ModeNotification:
		return nil, errors.New("notification sensor is not implemented")
	default:
		return nil, fmt.Errorf("%x is unknown sensor", mode)
	}
	return sensor.parse(header, h)
}

func validateData(data string) error {
	switch {
	case !strings.HasPrefix(data, ":"):
		return errors.New("data must have ':' as a prefix")
	case len(data) != 69:
		return errors.New("data size is invalid")
	}
	return nil
}

func verifyChecksum(data []byte) error {
	var sum byte
	for _, b := range data {
		sum += b
	}
	if sum != 0 {
		return errors.New("checksum is invalid")
	}
	return nil
}
