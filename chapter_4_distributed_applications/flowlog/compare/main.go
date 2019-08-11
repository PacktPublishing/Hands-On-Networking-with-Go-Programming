package main

import (
	"fmt"
	"net"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_4_distributed_applications/flowlog"
)

var history = []string{
	"2 123456789010 eni-abc123de 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1418530010 1418530070 ACCEPT OK",
	"2 123456789010 eni-abc123de 172.31.9.69 172.31.9.12 49761 3389 6 20 4249 1418530010 1418530070 REJECT OK",
	"2 123456789010 eni-1235b8ca 203.0.113.12 172.31.16.139 0 0 1 4 336 1432917027 1432917142 ACCEPT OK",
	"2 123456789010 eni-1235b8ca 172.31.16.139 203.0.113.12 0 0 1 4 336 1432917094 1432917142 REJECT OK",
}

func ipAddressAndPort(ip net.IP, port int64) string {
	return fmt.Sprintf("%v:%v", ip, port)
}

// MB is one megabyte of data.
const MB = 1024 * 1024

func isMoreThan1MBEgress(subnet *net.IPNet, r flowlog.Record) bool {
	return subnet.Contains(r.SourceAddress) && r.Bytes > MB
}

func neverConnectedBefore(subnet *net.IPNet, pastOutboundConnections map[string]struct{}, r flowlog.Record) bool {
	if !subnet.Contains(r.SourceAddress) {
		return false
	}
	_, connected := pastOutboundConnections[ipAddressAndPort(r.DestinationAddress, r.DestinationPort)]
	return !connected
}

func main() {
	// Get the old data.
	var records []flowlog.Record
	for _, l := range history {
		r, err := flowlog.Parse(l)
		if err != nil {
			fmt.Printf("Error parsing log: %v\n", err)
		}
		records = append(records, r)
	}

	// Create a history of past connections.
	pastOutboundConnections := make(map[string]struct{})
	for _, r := range records {
		pastOutboundConnections[ipAddressAndPort(r.DestinationAddress, r.DestinationPort)] = struct{}{}
	}

	// Define local subnet.
	_, subnet, err := net.ParseCIDR("172.31.0.0/16")
	if err != nil {
		fmt.Printf("Error parsing subnet: %v\n", err)
		return
	}

	// Parse the new data.
	newData := []string{
		// New IP address
		"2 123456789010 eni-abc123de 172.31.16.139 192.168.16.21 9999 443 6 20 4249 1418530010 1418540070 REJECT OK",
		// > 1MB egress
		"2 123456789010 eni-abc123de 172.31.16.139 192.168.16.21 9999 443 6 20 1048577 1418530010 1418540070 REJECT OK",
		// Nothing unusual
		"2 123456789010 eni-abc123de 172.31.16.139 192.168.16.21 9999 443 6 20 1024 1418530010 1418540070 REJECT OK",
	}

	for _, d := range newData {
		new, err := flowlog.Parse(d)
		if err != nil {
			fmt.Printf("Error parsing new log: %v\n", err)
			continue
		}
		if isMoreThan1MBEgress(subnet, new) {
			fmt.Printf("> 1MB egress from %v to %v\n", new.SourceAddress, ipAddressAndPort(new.DestinationAddress, new.DestinationPort))
		}
		if neverConnectedBefore(subnet, pastOutboundConnections, new) {
			fmt.Printf("New connection from %v to %v\n", new.SourceAddress, ipAddressAndPort(new.DestinationAddress, new.DestinationPort))
		}
		pastOutboundConnections[ipAddressAndPort(new.DestinationAddress, new.DestinationPort)] = struct{}{}
	}
}
