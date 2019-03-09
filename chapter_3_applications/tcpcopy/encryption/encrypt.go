package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func Encrypt(output io.Writer, input io.Reader, key []byte) (err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: input}
	_, err = io.Copy(output, reader)
	return
}

func Decrypt(output io.Writer, input io.Reader, key []byte) (err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{S: stream, W: output}
	_, err = io.Copy(writer, input)
	return
}
