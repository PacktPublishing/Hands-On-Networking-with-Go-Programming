package receive

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/tcpcopy/encryption"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/tcpcopy/wordlist"
)

// Start the receive.
func Start() (err error) {
	words, err := wordlist.Random(4)
	if err != nil {
		err = fmt.Errorf("unable to create encyrption key: %v", err)
		return
	}
	encryptionKey := sha256.Sum256([]byte(strings.Join(words, " ")))

	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 8002})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := l.AcceptTCP()
			if err != nil && err.Error() != "use of closed network connection" {
				fmt.Println("Error:", err)
				continue
			}
			go handle(conn, encryptionKey[:])
		}
	}()

	fmt.Printf("Listening on: %s\n", l.Addr())
	fmt.Println("Use command line:")
	fmt.Printf("  tcpcopy -m send -t %s -key '%s' -f filename.txt\n", l.Addr(), strings.Join(words, " "))
	<-sigs
	fmt.Println()
	fmt.Println("Closing down")
	l.Close()
	return
}

// ErrUnexpectedHash is returned when an unexpected hash is encountered.
var ErrUnexpectedHash = errors.New("unexpected hash, is the encryption key correct?")

var separator = []byte("|")

func handle(conn *net.TCPConn, encryptionKey []byte) {
	defer conn.Close()
	buf := make([]byte, 256*1024) // 256KB buffer.

	// Read the first 256KB.
	n, err := conn.Read(buf)
	if err != nil {
		//TODO: Push errors to a channel.
		fmt.Println(err)
		return
	}
	// Split until we get a pipe.
	// Should receive:
	// filename.txt|sha256_hash_of_encrypted_data|encrypted_filedata
	ranges := bytes.SplitN(buf[:n], separator, 3)
	if len(ranges) != 3 {
		//TODO: Push errors to a channel.
		fmt.Println("invalid data", len(ranges))
		return
	}
	_, fn := path.Split(string(ranges[0]))
	//TODO: Validate the filename.
	expectedHash := ranges[1]
	data := ranges[2]

	// Create the output file.
	if _, err := os.Stat(fn); err == nil {
		//TODO: Push errors to a channel.
		fmt.Println(os.ErrExist)
		return
	}
	f, err := os.Create(fn)
	if err != nil {
		//TODO: Push errors to a channel.
		fmt.Println(err)
		return
	}

	// Write the output data to a SHA256 calculation and the output file.
	actualHash := sha256.New()
	w := io.MultiWriter(f, actualHash)

	// Combine the data we've already read, and the rest of the data from the TCP connection.
	r := io.MultiReader(bytes.NewReader(data), conn)
	// Decrypt the body.
	err = encryption.Decrypt(w, r, encryptionKey)
	if err != nil {
		//TODO: Push errors to a channel.
		fmt.Println(err)
		return
	}

	// Check that the expected hash equals the data we got.
	if !areEqual(actualHash.Sum(nil), expectedHash) {
		err = ErrUnexpectedHash
	}
	return
}

func areEqual(a, b []byte) bool {
	if a != nil && b == nil || a == nil && b != nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
