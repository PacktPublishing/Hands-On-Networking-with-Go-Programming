package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/semaphore"
)

const maxConcurrentRequests int64 = 4

var accepted = []byte{1, 1, 1, 1, 1, 1, 1, 1}
var rejected = []byte{0, 0, 0, 0, 0, 0, 0, 0}

func main() {
	sem := semaphore.NewWeighted(maxConcurrentRequests)

	go func() {
		l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 8050})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for {
			conn, err := l.AcceptTCP()
			if err != nil && err.Error() != "use of closed network connection" {
				fmt.Println("Error:", err)
				continue
			}
			if canAcceptMore := sem.TryAcquire(1); canAcceptMore {
				go func() {
					defer conn.Close()
					defer sem.Release(1)
					fmt.Println("Accepted connection")
					if _, err := conn.Write(accepted); err != nil {
						fmt.Println("Error sending acceptance:", err)
						fmt.Println("Closing")
						return
					}
					data := make([]byte, 8)
					for {
						binary.LittleEndian.PutUint64(data, uint64(rand.Int63()))
						if _, err := conn.Write(data); err != nil {
							fmt.Println("Error sending data:", err)
							fmt.Println("Closing")
							return
						}
						time.Sleep(time.Second)
					}
				}()
			} else {
				fmt.Println("Rejected connection, max connections achieved")
				go func() {
					defer conn.Close()
					conn.Write(rejected)
				}()
			}
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
