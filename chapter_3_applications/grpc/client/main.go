package main

import (
	"context"
	"fmt"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/grpc/interface/order"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	client := order.NewOrdersClient(cc)
	addResp, err := client.Add(context.Background(), &order.AddRequest{
		Fields: &order.Fields{
			UserId: "userid123",
			Items: []*order.Item{
				&order.Item{
					Id:    "item1",
					Desc:  "description of the order",
					Price: 3200,
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(addResp)

	getResp, err := client.Get(context.Background(), &order.GetRequest{
		Id: addResp.Id,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(getResp)
}
