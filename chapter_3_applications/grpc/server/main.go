package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/grpc/interface/order"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc"
)

type OrderServer struct {
	orders map[string]*order.Order
}

func (os *OrderServer) Get(ctx context.Context, r *order.GetRequest) (resp *order.GetResponse, err error) {
	v, ok := os.orders[r.Id]
	if !ok {
		err = grpc.Errorf(codes.NotFound, "order %s not found", r.Id)
		return
	}
	resp = &order.GetResponse{
		Order: v,
	}
	return
}

func (os *OrderServer) Add(ctx context.Context, r *order.AddRequest) (resp *order.AddResponse, err error) {
	key := fmt.Sprintf("id%d", time.Now().UnixNano())
	os.orders[key] = &order.Order{
		Id:     key,
		Fields: r.Fields,
	}
	resp = &order.AddResponse{
		Id: key,
	}
	return
}

func main() {
	os := &OrderServer{
		orders: make(map[string]*order.Order),
	}
	s := grpc.NewServer()
	order.RegisterOrdersServer(s, os)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.Serve(l)
	if err != nil {
		fmt.Println(err)
		return
	}
}
