package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	sc, err := ioutil.ReadFile("../cert.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(sc))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: roots,
		},
	}
	c := &http.Client{Transport: tr}
	resp, err := c.Get("https://adrian.local:8443")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
