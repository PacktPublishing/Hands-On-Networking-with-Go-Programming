package main

import (
	"fmt"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/client"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

func main() {
	cc := client.NewCompanyClient("http://localhost:9021")

	fmt.Println("Listing companies")
	companies, err := cc.List()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range companies {
		fmt.Println(c)
	}
	fmt.Println()

	fmt.Println("Creating new company")
	c := company.Company{
		Name: "Test Company 1",
		Address: company.Address{
			Address1: "Address 1",
			Address2: "Address 2",
			Address3: "Address 3",
			Address4: "Address 4",
			Postcode: "Postcode",
		},
		VATNumber: "12345",
	}
	id, err := cc.Post(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Created ID: %d\n", id)
	fmt.Println()

	fmt.Println("Getting company")
	c, err = cc.Get(id.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()

	fmt.Println("Listing companies")
	companies, err = cc.List()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range companies {
		fmt.Println(c)
	}
	fmt.Println()
}
