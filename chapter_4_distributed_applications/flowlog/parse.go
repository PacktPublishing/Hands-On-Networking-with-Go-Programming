package flowlog

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const versionFieldIndex = 0
const accountIDFieldIndex = 1
const interfaceIDFieldIndex = 2
const srcAddrFieldIndex = 3
const dstAddrFieldIndex = 4
const srcPortFieldIndex = 5
const dstPortFieldIndex = 6
const protocolFieldIndex = 7
const packetsFieldIndex = 8
const bytesFieldIndex = 9
const startFieldIndex = 10
const endFieldIndex = 11
const actionFieldIndex = 12
const logStatusFieldIndex = 13

const fieldCount = 14

// Parse the AWS flow log Record.
func Parse(s string) (Record, error) {
	split := strings.SplitN(s, " ", fieldCount)
	if len(split) != fieldCount {
		return Record{}, fmt.Errorf("flowlog: invalid field count: %v", len(split))
	}
	startSeconds, err := strconv.ParseInt(split[startFieldIndex], 10, 64)
	if err != nil {
		return Record{}, fmt.Errorf("flowlog: invalid start time: '%v'", split[startFieldIndex])
	}
	endSeconds, err := strconv.ParseInt(split[endFieldIndex], 10, 64)
	if err != nil {
		return Record{}, fmt.Errorf("flowlog: invalid end time: '%v'", split[endFieldIndex])
	}

	return Record{
		Version:            split[versionFieldIndex],
		AccountID:          split[accountIDFieldIndex],
		InterfaceID:        split[interfaceIDFieldIndex],
		SourceAddress:      net.ParseIP(split[srcAddrFieldIndex]),
		DestinationAddress: net.ParseIP(split[dstAddrFieldIndex]),
		SourcePort:         parseIntOrZero(split[srcPortFieldIndex]),
		DestinationPort:    parseIntOrZero(split[dstPortFieldIndex]),
		Protocol:           split[protocolFieldIndex],
		Packets:            parseIntOrZero(split[packetsFieldIndex]),
		Bytes:              parseIntOrZero(split[bytesFieldIndex]),
		Start:              time.Unix(startSeconds, 0),
		End:                time.Unix(endSeconds, 0),
		Action:             Action(split[actionFieldIndex]),
		Status:             Status(split[logStatusFieldIndex]),
	}, nil
}

func parseIntOrZero(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return i
}

// Action of the Record.
type Action string

const (
	// ActionNone is the action when no data is present.
	ActionNone Action = "-"
	// ActionAccept is the action when it the connection was accepted.
	ActionAccept Action = "ACCEPT"
	// ActionReject is the action when it the connection was rejected.
	ActionReject Action = "REJECT"
)

// Status of the Record.
type Status string

const (
	// StatusOK is a successful log.
	StatusOK Status = "OK"
	// StatusNoData is when the log has no data in the collection period.
	StatusNoData Status = "NODATA"
	// StatusSkipData is the status when records were skipped in the capture window.
	StatusSkipData Status = "SKIPDATA"
)

// Record represents a network flow in your flow log.
type Record struct {
	// Version of the VPC Flow Logs.
	Version string
	// AccountID for the flow log.
	AccountID string
	// InterfaceID of the network interface for which the traffic is recorded.
	InterfaceID string
	// SourceAddress is an IPv4 or IPv6 address. The IPv4 address of the network interface is always its private IPv4 address.
	SourceAddress net.IP
	// DestinationAddress is an IPv4 or IPv6 address. The IPv4 address of the network interface is always its private IPv4 address.
	DestinationAddress net.IP
	// SourcePort of the traffic.
	SourcePort int64
	// DestinationPort of the traffic.
	DestinationPort int64
	// Protocol is the IANA protocol number of the traffic. For more information, see Assigned Internet Protocol Numbers.
	Protocol string
	// 	Packets is the number of packets transferred during the capture window.
	Packets int64
	// Bytes is the number of bytes transferred during the capture window.
	Bytes int64
	// Start time of the start of the capture window.
	Start time.Time
	// End time of the end of the capture window.
	End time.Time
	// Action associated with the traffic - ACCEPT: The recorded traffic was permitted by the security groups or network ACLs. or REJECT: The recorded traffic was not permitted by the security groups or network ACLs.
	Action Action
	// Status of the flow log: OK: Data is logging normally to the chosen destinations. NODATA: There was no network traffic to or from the network interface during the capture window. SKIPDATA: Some flow log records were skipped during the capture window. This may be because of an internal capacity constraint, or an internal error.
	Status Status
}
