package main

import (
	"fmt"
	"time"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_4_distributed_applications/stream"
)

func main() {
	exit := func(head stream.Event, tail []stream.Event) {
		if head.Type != "itemAddedToBasket" {
			return
		}
		var success bool
		for _, t := range tail {
			if t.SessionID == head.SessionID && t.Type == "paymentReceived" {
				success = true
			}
		}
		fmt.Printf("%s: succeeded: %v\n", head.SessionID, success)
	}
	w := stream.NewEventWindow(time.Hour*5, exit)

	// Simulate events.
	events := []stream.Event{
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Type:      "itemViewed",
			SessionID: "SessionA",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 0, 10, 0, 0, time.UTC),
			Type:      "itemViewed",
			SessionID: "SessionB",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC),
			Type:      "itemAddedToBasket",
			SessionID: "SessionA",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
			Type:      "itemAddedToBasket",
			SessionID: "SessionB",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 2, 0, 0, 0, time.UTC),
			Type:      "checkoutScreenViewed",
			SessionID: "SessionA",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 3, 0, 0, 0, time.UTC),
			Type:      "paymentScreenViewed",
			SessionID: "SessionA",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 1, 4, 0, 0, 0, time.UTC),
			Type:      "paymentReceived",
			SessionID: "SessionA",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			Type:      "itemViewed",
			SessionID: "SessionC",
		},
		stream.Event{
			Date:      time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			Type:      "itemViewed",
			SessionID: "SessionC",
		},
	}
	for _, e := range events {
		w.Push(e)
	}
}
