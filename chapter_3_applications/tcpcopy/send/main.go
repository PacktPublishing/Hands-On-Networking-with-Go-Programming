package send

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/tcpcopy/encryption"
)

func Start(fileName string, url string, key string) (err error) {
	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	encryptionKey := sha256.Sum256([]byte(key))

	hash := sha256.New()
	err = encryption.Encrypt(hash, file, encryptionKey[:])
	if err != nil {
		return
	}
	// Reset the file, because we need to read it again.
	_, err = file.Seek(0, 0)
	if err != nil {
		return
	}

	conn, err := net.Dial("tcp", url)
	if err != nil {
		return
	}
	defer conn.Close()
	// Write the header.
	_, fn := path.Split(fileName)
	header := fn + "|" + fmt.Sprintf("%2x", hash.Sum(nil)) + "|"
	_, err = conn.Write([]byte(header))
	if err != nil {
		err = fmt.Errorf("failed to write header: %v", err)
		return
	}
	// Write the rest of the file.
	err = encryption.Encrypt(conn, file, encryptionKey[:])
	return
}
