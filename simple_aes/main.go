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

func encrypt(cli_key []byte, iv []byte, src string, dest string) {

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write(cli_key)
	key := sum.Sum(nil)

	// build the iv with sha1 hash
	ivhash := md5.New()
	ivhash.Write(iv)
	iv_arr := ivhash.Sum(nil)

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
	if _, err := io.ReadFull(rand.Reader, iv_arr); err != nil {
		log.Fatalln("error while reading iv")
	}

	// encrypt the message
	stream := cipher.NewCBCEncrypter(block, iv_arr)
	stream.CryptBlocks(enc_message[aes.BlockSize:], f)

	fmt.Println("%x00", enc_message)
}

func decrypt(cli_key []byte, iv []byte, src string) {

	// build the key with sha256 hash
	sum := sha256.New()
	sum.Write(cli_key)
	key := sum.Sum(nil)

	// build the iv with sha1 hash
	ivhash := md5.New()
	ivhash.Write(iv)
	iv_arr := ivhash.Sum(nil)

	f, err := os.ReadFile(src)
	if err != nil {
		log.Fatalln("error while reading src")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("error while building aes key")
	}

	decrypted_message := f[aes.BlockSize:]

	// decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv_arr)
	stream.CryptBlocks(decrypted_message, decrypted_message)

	fmt.Println("Encrypted Message...")
	fmt.Println(string(decrypted_message))
}

func main() {

	args := os.Args

	if len(args) < 2 || len(args) > 6 {
		fmt.Println("./simple_aes mode[encrypt:decrypt] key|'' iv|'' src|'' dest|''")
		return
	}

	switch args[1] {
	case "e":
		encrypt([]byte(args[2]), []byte(args[3]), args[4], args[5])
	case "d":
		decrypt([]byte(args[2]), []byte(args[3]), args[4])
	}

}
