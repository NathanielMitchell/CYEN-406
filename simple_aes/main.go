package simple_aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Inputs:
// key is a hashed key sha256
// iv is a hashed md5 value

func encrypt(key []byte, iv []byte, src string, dest string) {

	f, err := os.ReadFile(src)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println("error while building aes key")
	}

	enc_message := make([]byte, aes.BlockSize+len(f))
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("error while reading iv")
	}

	// encrypt the message
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(enc_message[aes.BlockSize:], f)

	fmt.Printf("Writing hex-encoded message to %s...\n", dest)
	hex_enc_message := hex.EncodeToString(enc_message)
	os.WriteFile(dest, []byte(hex_enc_message), 0640)
}

// key is a sha256 hash
// iv is a md5 hash
// src is a hex-encoded ciphertext
func decrypt(key []byte, iv []byte, src string, dest string) {

	fcontents, err := os.ReadFile(src)
	if err != nil {
		fmt.Println(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error while building aes key")
	}

	decrypted_message, err := hex.DecodeString(string(fcontents))
	if err != nil {
		fmt.Println(err)
	}
	decrypted_message = decrypted_message[aes.BlockSize:]

	// decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(decrypted_message, decrypted_message)

	fmt.Printf("Encrypted Message written to %s...\n", dest)
	os.WriteFile(dest, decrypted_message, 0640)
}

func simple_aes() {

	args := os.Args

	if len(args) < 2 || len(args) > 6 {
		fmt.Println("./simple_aes mode[e:d] key iv src dest")
		return
	}

	dbarray, err := hex.DecodeString(args[2])
	if err != nil {
		fmt.Print("error")
	}

	harray, err := hex.DecodeString(args[3])
	if err != nil {
		fmt.Print("error")
	}

	switch args[1] {
	case "e":
		encrypt(dbarray, harray, args[4], args[5])
	case "d":
		decrypt(dbarray, harray, args[4], args[5])
	}

}
