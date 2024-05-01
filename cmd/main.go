package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/NathanielMitchell/CYEN-406/dhke"
	"github.com/NathanielMitchell/CYEN-406/simple_aes"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "9001"
	SERVER_TYPE = "tcp"
)

func main() {

	args := os.Args

	// username/password combo
	// this should make it single simple string for dhke
	combo := args[1]

    server(combo)

	for true {
		// message to AES encrypt
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("message to send: ")
		message, _ := reader.ReadString('\n')

		message = simple_aes.Encrypt([]byte(message), resultGuy.symkey, resultGuy.iv)
        con := processClient()
	}
}
