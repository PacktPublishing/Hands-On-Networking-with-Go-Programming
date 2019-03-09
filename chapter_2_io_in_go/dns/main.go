package main

import (
	"context"
	"flag"
	"fmt"
	"net"
)

var hostFlag = flag.String("host", "google.com", "The DNS name of the host to lookup IP addresses for")

func main() {
	flag.Parse()
	host := *hostFlag

	var r net.Resolver
	r.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		fmt.Println("was going to use DNS server at", address)
		// Google's is 8.8.8.8
		// OpenDNS at 208.67.222.123
		return net.Dial(network, "208.67.222.123:53")
	}
	r.PreferGo = true
	res, err := r.LookupHost(context.Background(), host)
	fmt.Println("IP Addresses", res, err)
}
