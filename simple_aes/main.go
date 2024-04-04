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

func main() {

	args := os.Args

	if len(args) != 4 {
		fmt.Println("./simple_aes key src dest")
		return
	}

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write([]byte(args[1]))
	key := sum.Sum(nil)

	f, err := os.ReadFile(args[2])
	if err != nil {
		log.Fatalln("error while reading src")
	}

	// pad the end of a block
	if len(f)%aes.BlockSize != 0 {
		missingBytes := len(f) % aes.BlockSize
		totalPad := aes.BlockSize - missingBytes
		for i := 0; i < totalPad; i++ {
			// the code I found that described how to do this wanted the actual pad value to be the same as 'totalPad'
			f = append(f, byte(0x00))
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("error while building aes key")
	}

	enc_message := make([]byte, aes.BlockSize+len(f))
	iv := enc_message[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalln("error while reading iv")
	}

	// encrypt the message
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(enc_message[aes.BlockSize:], f)

	os.WriteFile(args[3], enc_message, 0644)

}
