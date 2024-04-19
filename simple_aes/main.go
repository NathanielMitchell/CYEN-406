package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
    "crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func encrypt(cli_key []byte, iv []byte, src []byte) {

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write(cli_key)
	key := sum.Sum(nil)

	// build the iv with sha1 hash
	ivhash := md5.New()
	ivhash.Write(iv)
	iv_arr := ivhash.Sum(nil)

	// pad the end of a block
	if len(src)%aes.BlockSize != 0 {
		missingBytes := len(src) % aes.BlockSize
		totalPad := aes.BlockSize - missingBytes
		for i := 0; i < totalPad; i++ {
			// the code that described how to do this
			// wanted the actual pad value to be the same as 'totalPad'
			src = append(src, byte(0x00))
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("error while building aes key")
	}

	enc_message := make([]byte, aes.BlockSize+len(src))
	if _, err := io.ReadFull(rand.Reader, iv_arr); err != nil {
		log.Fatalln("error while reading iv")
	}

	// encrypt the message
	stream := cipher.NewCBCEncrypter(block, iv_arr)
	stream.CryptBlocks(enc_message[aes.BlockSize:], src)

	fmt.Println("%x00", enc_message)
}

func decrypt(cli_key []byte, iv []byte, src []byte) {

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write(cli_key)
	key := sum.Sum(nil)

	// build the iv with sha1 hash
	ivhash := md5.New()
	ivhash.Write(iv)
	iv_arr := ivhash.Sum(nil)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("error while building aes key")
	}

	decrypted_message := src[aes.BlockSize:]

	// decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv_arr)
	stream.CryptBlocks(decrypted_message, decrypted_message)

	fmt.Println("Encrypted Message...")
	fmt.Println(string(decrypted_message))
}

func main() {

	args := os.Args

	if len(args) < 2 || len(args) > 6 {
		fmt.Println("./simple_aes mode[encrypt:decrypt] key|'' iv|''")
		return
	}

	switch args[1] {
	case "e":
		encrypt([]byte(args[2]), []byte(args[3]), []byte(args[4]))
	case "d":
		decrypt([]byte(args[2]), []byte(args[3]), []byte(args[4]))
	}

}
