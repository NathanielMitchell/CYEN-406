package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func checkErr(err error) {
	log.Fatal(err)
}

func main() {

	args := os.Args

	if len(args) != 4 {
		fmt.Println("./simple_aes key src dest")
		return
	}

	// build the key with sha256sum
	sum := sha256.New()
	key := sum.Sum([]byte(args[1]))

	f, err := os.ReadFile(args[2])
	checkErr(err)

	block, err := aes.NewCipher(key)
	checkErr(err)

	// encrypt the message
	enc_message := make([]byte, aes.BlockSize+len(f))
	iv := enc_message[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(enc_message[aes.BlockSize:], f)

	os.WriteFile(args[3], enc_message, 0644)

}
