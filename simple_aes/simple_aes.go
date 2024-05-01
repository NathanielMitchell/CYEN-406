package simple_aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Inputs:
// key is a hashed key sha256
// iv is a hashed md5 value

func Encrypt(key []byte, iv []byte, message []byte) *[]byte {

	// pad the end of a block
	if len(message)%aes.BlockSize != 0 {
		missingBytes := len(message) % aes.BlockSize
		totalPad := aes.BlockSize - missingBytes
		for i := 0; i < totalPad; i++ {
			// the code that described how to do this
			// wanted the actual pad value to be the same as 'totalPad'
			message = append(message, byte(0x00))
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error while building aes key")
	}

	enc_message := make([]byte, aes.BlockSize+len(message))
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("error while reading iv")
	}

	// encrypt the message
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(enc_message[aes.BlockSize:], message)

    return &enc_message
}

// key is a sha256 hash
// iv is a md5 hash
// src is a hex-encoded ciphertext
func Decrypt(key []byte, iv []byte, enc_message []byte) *[]byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error while building aes key")
	}

	decrypted_message, err := hex.DecodeString(string(enc_message))
	if err != nil {
		fmt.Println(err)
	}
	decrypted_message = decrypted_message[aes.BlockSize:]

	// decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(decrypted_message, decrypted_message)

    return &decrypted_message
}
