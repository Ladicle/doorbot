package sensor

import (
	"errors"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		wantSensor CloserSensor
		wantError  error
	}{
		{
			name:      "empty data",
			wantError: errors.New("data must have ':' as a prefix"),
		},
		{
			name:      "invalid length",
			data:      ":0123",
			wantError: errors.New("data size is invalid"),
		},
		{
			name:      "invalid checksum",
			data:      ":80000000A8001C82012B1E01808103113008020D0C1130010203E40000000101EC6A",
			wantError: errors.New("checksum is invalid"),
		},
		{
			name:      "non closer sensor",
			data:      ":80000000A8001C82012B1E01808203113008020D0C1130010203E40000000101EC6D",
			wantError: errors.New("environment sensor is not implemented"),
		},
		{
			name: "normal",
			data: ":80000000A8001C82012B1E01808103113008020D0C1130010203E40000000101EC6E",
			wantSensor: CloserSensor{
				commonHeader: commonHeader{
					SerialID:     "82012B1E",
					DeviceID:     "01",
					LQI:          168,
					PowerVoltage: 3340,
				},
				State: CloseNPole,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := Parse(tt.data)
			if err != nil {
				if tt.wantError == nil || tt.wantError.Error() != err.Error() {
					t.Fatalf("unexpected error occurred: %v", err)
				}
				return
			}
			if tt.wantError != nil {
				t.Fatalf("expected %q error but no error occurred", tt.wantError)
			}
			if !reflect.DeepEqual(cs, tt.wantSensor) {
				t.Errorf("CloseSensor has unexpected values:\nwant=%+v\ngot=%+v",
					cs, tt.wantSensor)
			}
		})
	}
}
