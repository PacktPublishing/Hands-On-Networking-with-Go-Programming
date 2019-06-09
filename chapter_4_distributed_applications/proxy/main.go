package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/armon/go-socks5"
)

func main() {
	// Start up an server to warn about disallowed domains.
	go func() {
		http.ListenAndServe(":10000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Accessing warning page.")
			w.Write([]byte("The domain you tried to access is not permitted."))
		}))
	}()

	go func() {
		// Generate your own certificate against your local CA using the steps at https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/
		certFile := "/Users/adrian/Documents/localca/LocalCAAll.crt"
		keyFile := "/Users/adrian/Documents/localca/LocalCAAll.key"
		http.ListenAndServeTLS(":10001", certFile, keyFile, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Accessing TLS warning page.")
			w.Write([]byte("The TLS domain you tried to access is not permitted."))
		}))
	}()

	// Redirect all disallowed IP addresses to our local server.
	s, err := socks5.New(&socks5.Config{
		Rewriter: rewriter{},
	})
	if err != nil {
		log.Fatal(err)
	}
	s.ListenAndServe("tcp", ":9999")
}

type rewriter struct{}

func (ur rewriter) Rewrite(ctx context.Context, r *socks5.Request) (context.Context, *socks5.AddrSpec) {
	fmt.Printf("Atttempting to connect to: %v\n", r.DestAddr.IP.String())
	allowed, err := isAllowedIP(r.DestAddr.IP)
	if err != nil {
		log.Printf("failed to get IPs: %v", err)
	}
	if !allowed {
		// Rewrite to use our local server.
		fmt.Println("Rewriting to use local server.")
		if r.DestAddr.Port == 80 || r.DestAddr.Port == 8080 {
			// Connect to our local insecure Web server.
			return ctx, &socks5.AddrSpec{
				FQDN: "localhost",
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 10000,
			}
		}
		// Connect to the TLS endpoint.
		return ctx, &socks5.AddrSpec{
			FQDN: "localhost",
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 10001,
		}
	}
	fmt.Printf("Done.")
	return ctx, r.DestAddr
}

func getDisallowedIPAddresses() (disallowed []net.IP, err error) {
	disallowedDomains := []string{"google.com"}
	var ips []net.IP
	for _, disallowedDomain := range disallowedDomains {
		ips, err = net.LookupIP(disallowedDomain)
		if err != nil {
			return
		}
		disallowed = append(disallowed, ips...)
	}
	fmt.Printf("Disallowed domains: %v\n", disallowed)
	return
}

func isAllowedIP(ip net.IP) (bool, error) {
	disallowed, err := getDisallowedIPAddresses()
	if err != nil {
		return false, err
	}
	for _, d := range disallowed {
		fmt.Printf("Checking: %v is in disallowed list of %v\n", ip, disallowed)
		if ip.Equal(d) {
			return false, nil
		}
	}
	return true, nil
}
