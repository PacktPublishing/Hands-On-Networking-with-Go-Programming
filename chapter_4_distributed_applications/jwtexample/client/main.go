package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	username := "user1@example.com"
	pwd := "BadPassword1"
	r, err := http.NewRequest("GET", "http://localhost:8000", nil)
	if err != nil {
		fmt.Println("error create request to issuer:", err)
		os.Exit(1)
	}
	r.SetBasicAuth(username, pwd)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println("error getting token from issuer:", err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpeted status code from issuer:", resp.StatusCode)
		os.Exit(1)
	}
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading token from issuer:", err)
		os.Exit(1)
	}
	fmt.Println("Got token from issuer. Making request to server.")

	// Now lets make a request to the server.
	sr, err := http.NewRequest("GET", "http://localhost:8001/whoami", nil)
	if err != nil {
		fmt.Println("error creating request to server:", err)
		os.Exit(1)
	}
	sr.Header.Set("Authorization", "Bearer "+string(token))
	resp, err = http.DefaultClient.Do(sr)
	if err != nil {
		fmt.Println("error getting response from server:", err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpeted status code from server:", resp.StatusCode)
		os.Exit(1)
	}
	whoami, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading whoami response from server:", err)
		os.Exit(1)
	}
	fmt.Println("Made request to server.")
	fmt.Println("whoami:", string(whoami))
}
