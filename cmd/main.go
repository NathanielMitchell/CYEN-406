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

	// ip for other team
    // means that we need to start dh
	ip := args[2]
    if len(ip) != 0 {

    }

	con := connectionHandler{[]byte{}, []byte{}, nil}
	go server(con)

	key, iv, err := dhke.ServerDhke(combo, con)
	if err != nil {
		fmt.Println("error while trying to run dh key exchange")
	}

	for true {
		// message to AES encrypt
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("message to send: ")
		message, _ := reader.ReadString('\n')

		message = simple_aes.Encrypt([]byte(message), key, iv)
		client([]byte(message), ip)
	}
}
