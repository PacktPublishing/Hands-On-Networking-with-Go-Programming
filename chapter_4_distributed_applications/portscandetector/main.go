package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// NewCappedCounter creates a counter for string values which has a maximum
// number of values it can contain.
func NewCappedCounter() *CappedCounter {
	return &CappedCounter{
		Capacity: 100,
		Keys:     []string{},
		Values:   make(map[string]int),
	}
}

// CappedCounter is a counter of how many times a key has been seen
// but caps the maximum number of keys it retains.
type CappedCounter struct {
	Capacity int
	Keys     []string
	Values   map[string]int
	m        sync.Mutex
}

// Increment the key's value by one and return the total.
func (cc *CappedCounter) Increment(k string) (total int) {
	cc.m.Lock()
	defer cc.m.Unlock()
	v, ok := cc.Values[k]
	if !ok {
		cc.Keys = append(cc.Keys, k)
		if len(cc.Keys) > cc.Capacity {
			cc.Keys = cc.Keys[1:]
		}
	}
	total = v + 1
	cc.Values[k] = total
	return
}

// Scan is the source and target of the scan.
type Scan struct {
	From string
	To   string
}

// Start a Listener.
func Start(port int, scanned chan Scan) (stop func(), err error) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	stop = func() {
		wg.Wait()
		l.Close()
	}
	go func(l *net.TCPListener, scanned chan Scan) {
		defer wg.Done()
		for {
			conn, err := l.AcceptTCP()
			if err != nil {
				if err.Error() == "use of closed network connection" {
					return
				}
				continue
			}
			go func() {
				defer conn.Close()
				segs := strings.Split(conn.RemoteAddr().String(), ":")
				fromIP := strings.Join(segs[:len(segs)-1], ":")
				scanned <- Scan{From: fromIP, To: conn.LocalAddr().String()}
			}()
		}
	}(l, scanned)
	return
}

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)

	ports := []int{8080, 8086, 9100}
	stoppers := make([]func(), len(ports))
	scanned := make(chan Scan)
	for i, p := range ports {
		var err error
		stoppers[i], err = Start(p, scanned)
		if err != nil {
			fmt.Println(err)
		}
	}
	count := NewCappedCounter()

	for {
		select {
		case by := <-scanned:
			total := count.Increment(by.From)
			fmt.Printf("Connection from %s to %s\n", by.From, by.To)
			if total >= 3 {
				fmt.Printf("Scanned by %s\n", by.From)
			}
		case <-sigs:
			fmt.Println("Shutting down...")
			for _, stop := range stoppers {
				stop()
			}
			return
		}
	}
}
