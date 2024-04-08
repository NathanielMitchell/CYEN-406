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

func encrypt(cli_key []byte, src string, dest string) {

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write(cli_key)
	key := sum.Sum(nil)

	f, err := os.ReadFile(src)
	if err != nil {
		log.Fatalln("error while reading src")
	}

	// pad the end of a block
	if len(f)%aes.BlockSize != 0 {
		missingBytes := len(f) % aes.BlockSize
		totalPad := aes.BlockSize - missingBytes
		for i := 0; i < totalPad; i++ {
			// the code that described how to do this
			// wanted the actual pad value to be the same as 'totalPad'
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

	os.WriteFile(dest, enc_message, 0644)
}

func decrypt(cli_key []byte, src string) {

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write(cli_key)
	key := sum.Sum(nil)

	f, err := os.ReadFile(src)
	if err != nil {
		log.Fatalln("error while reading src")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("error while building aes key")
	}

	decrypted_message := f[aes.BlockSize:]
	iv := decrypted_message[:aes.BlockSize]

	// decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(decrypted_message, decrypted_message)

	fmt.Println("Encrypted Message...")
	fmt.Println(decrypted_message)
}

func main() {

	args := os.Args

	if len(args) != 4 {
		fmt.Println("./simple_aes mode[e:d] ke|'' iv|'' src|'' dest|''")
		return
	}

	switch args[1] {
	case "e":
		encrypt([]byte(args[2]), args[3], args[4])
	case "d":
		decrypt([]byte(args[2]), args[3])
	}

}
