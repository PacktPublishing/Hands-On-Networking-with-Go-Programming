package flowlog

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Record
	}{
		{
			name:  "SSH traffic (destination port 22, TCP protocol) to network interface eni-abc123de in account 123456789010 was allowed",
			input: "2 123456789010 eni-abc123de 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1418530010 1418530070 ACCEPT OK",
			expected: Record{
				Version:            "2",
				AccountID:          "123456789010",
				InterfaceID:        "eni-abc123de",
				SourceAddress:      net.ParseIP("172.31.16.139"),
				DestinationAddress: net.ParseIP("172.31.16.21"),
				SourcePort:         20641,
				DestinationPort:    22,
				Protocol:           "6",
				Packets:            20,
				Bytes:              4249,
				Start:              time.Unix(1418530010, 0),
				End:                time.Unix(1418530070, 0),
				Action:             ActionAccept,
				Status:             StatusOK,
			},
		},
		{
			name:  "RDP traffic (destination port 3389, TCP protocol) to network interface eni-abc123de in account 123456789010 was rejected",
			input: "2 123456789010 eni-abc123de 172.31.9.69 172.31.9.12 49761 3389 6 20 4249 1418530010 1418530070 REJECT OK",
			expected: Record{
				Version:            "2",
				AccountID:          "123456789010",
				InterfaceID:        "eni-abc123de",
				SourceAddress:      net.ParseIP("172.31.9.69"),
				DestinationAddress: net.ParseIP("172.31.9.12"),
				SourcePort:         49761,
				DestinationPort:    3389,
				Protocol:           "6",
				Packets:            20,
				Bytes:              4249,
				Start:              time.Unix(1418530010, 0),
				End:                time.Unix(1418530070, 0),
				Action:             ActionReject,
				Status:             StatusOK,
			},
		},
		{
			name:  "no data in window",
			input: "2 123456789010 eni-1a2b3c4d - - - - - - - 1431280876 1431280934 - NODATA",
			expected: Record{
				Version:            "2",
				AccountID:          "123456789010",
				InterfaceID:        "eni-1a2b3c4d",
				SourceAddress:      nil,
				DestinationAddress: nil,
				SourcePort:         0,
				DestinationPort:    0,
				Protocol:           "-",
				Packets:            0,
				Bytes:              0,
				Start:              time.Unix(1431280876, 0),
				End:                time.Unix(1431280934, 0),
				Action:             ActionNone,
				Status:             StatusNoData,
			},
		},
		{
			name:  "data in window skipped",
			input: "2 123456789010 eni-4b118871 - - - - - - - 1431280876 1431280934 - SKIPDATA",
			expected: Record{
				Version:            "2",
				AccountID:          "123456789010",
				InterfaceID:        "eni-4b118871",
				SourceAddress:      nil,
				DestinationAddress: nil,
				SourcePort:         0,
				DestinationPort:    0,
				Protocol:           "-",
				Packets:            0,
				Bytes:              0,
				Start:              time.Unix(1431280876, 0),
				End:                time.Unix(1431280934, 0),
				Action:             ActionNone,
				Status:             StatusSkipData,
			},
		},
		{
			name:  "ping request accept",
			input: "2 123456789010 eni-1235b8ca 203.0.113.12 172.31.16.139 0 0 1 4 336 1432917027 1432917142 ACCEPT OK",
			expected: Record{
				Version:            "2",
				AccountID:          "123456789010",
				InterfaceID:        "eni-1235b8ca",
				SourceAddress:      net.ParseIP("203.0.113.12"),
				DestinationAddress: net.ParseIP("172.31.16.139"),
				SourcePort:         0,
				DestinationPort:    0,
				Protocol:           "1",
				Packets:            4,
				Bytes:              336,
				Start:              time.Unix(1432917027, 0),
				End:                time.Unix(1432917142, 0),
				Action:             ActionAccept,
				Status:             StatusOK,
			},
		},
		{
			name:  "ping response reject",
			input: "2 123456789010 eni-1235b8ca 172.31.16.139 203.0.113.12 0 0 1 4 336 1432917094 1432917142 REJECT OK",
			expected: Record{
				Version:            "2",
				AccountID:          "123456789010",
				InterfaceID:        "eni-1235b8ca",
				SourceAddress:      net.ParseIP("172.31.16.139"),
				DestinationAddress: net.ParseIP("203.0.113.12"),
				SourcePort:         0,
				DestinationPort:    0,
				Protocol:           "1",
				Packets:            4,
				Bytes:              336,
				Start:              time.Unix(1432917094, 0),
				End:                time.Unix(1432917142, 0),
				Action:             ActionReject,
				Status:             StatusOK,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual, err := Parse(test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(test.expected, actual) {
				t.Errorf("expected:\n%+v\ngot:\n%+v", test.expected, actual)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "empty",
			input:    "",
			expected: fmt.Errorf("flowlog: invalid field count: 1"),
		},
		{
			name:     "invalid start time",
			input:    "2 123456789010 eni-abc123de 172.31.9.69 172.31.9.12 49761 3389 6 20 4249 not_time 1418530070 REJECT OK",
			expected: fmt.Errorf("flowlog: invalid start time: 'not_time'"),
		},
		{
			name:     "invalid end time",
			input:    "2 123456789010 eni-abc123de 172.31.9.69 172.31.9.12 49761 3389 6 20 4249 1418530070 not_time REJECT OK",
			expected: fmt.Errorf("flowlog: invalid end time: 'not_time'"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := Parse(test.input)
			if err == nil {
				t.Error("unexpected success")
				return
			}
			if test.expected.Error() != err.Error() {
				t.Errorf("expected:\n%s\ngot:\n%s", test.expected, err.Error())
			}
		})
	}
}
