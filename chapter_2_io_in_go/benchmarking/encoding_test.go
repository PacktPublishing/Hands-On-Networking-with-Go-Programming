package benchmarking

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"
	"time"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_2_io_in_go/benchmarking/messages"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type order struct {
	ID        int64     `json:"id"`
	Customer  customer  `json:"customer"`
	Items     []item    `json:"items"`
	OrderDate time.Time `json:"date"`
}

type customer struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type item struct {
	ID   int64  `json:"id"`
	URL  string `json:"url"`
	Desc string `json:"desc"`
	Cost int64  `json:"cost"`
}

var testOrder = order{
	ID: 1,
	Customer: customer{
		ID:        2,
		Email:     "customer@example.com",
		FirstName: "Joe",
		LastName:  "Tegmark",
	},
	Items: []item{
		item{
			ID:   3,
			URL:  "https://example.com/products/3",
			Desc: "item 3",
			Cost: 1553,
		},
		item{
			ID:   4,
			URL:  "https://example.com/products/4",
			Desc: "item 4",
			Cost: 12378,
		},
		item{
			ID:   5,
			URL:  "https://example.com/products/5",
			Desc: "item 5",
			Cost: 112,
		},
	},
	OrderDate: time.Now(),
}

func BenchmarkJSONEncoding(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(testOrder); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGobEncoding(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(testOrder); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkProtobufEncoding(b *testing.B) {
	getItems := func(items []item) (op []*messages.Item) {
		for _, itm := range items {
			op = append(op, &messages.Item{
				ItemId:      itm.ID,
				Url:         itm.URL,
				Description: itm.Desc,
				Cost:        itm.Cost,
			})
		}
		return
	}
	ts, err := ptypes.TimestampProto(testOrder.OrderDate)
	if err != nil {
		b.Fatal(err)
	}
	m := messages.Order{
		OrderId: testOrder.ID,
		Customer: &messages.Customer{
			CustomerId: testOrder.Customer.ID,
			FirstName:  testOrder.Customer.FirstName,
			LastName:   testOrder.Customer.LastName,
			Email:      testOrder.Customer.Email,
		},
		Items:     getItems(testOrder.Items),
		OrderDate: ts,
	}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		data, err := proto.Marshal(&m)
		if err != nil {
			b.Error(err)
		}
		if len(data) == 0 {
			b.Error("no bytes written")
		}
	}
}
