package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "9001"
	SERVER_TYPE = "tcp"
)

func hex_encode()

func main() {

	args := os.Args

	// flag for client/server usage
	// should be client|server
	mode := args[1]

	// username/password combo
	// this should make it single simple string for dhke
	combo := args[2]

	// ip for other team
	ip := args[3]

	for true {
		// message to AES encrypt
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("message to send: ")
		message, _ := reader.ReadString('\n')

		// run the program based on the mode supplied
		switch mode {
		case "server":
			dhke.dhke(combo)
			server()
		case "client":
			message := simple_aes.encrypt([]byte(message))
			client(message, ip)
		}

		dbarray, err := hex.DecodeString(args[2])
		if err != nil {
			fmt.Print("error while decoding the hex key")
		}

		harray, err := hex.DecodeString(args[3])
		if err != nil {
			fmt.Print("error while decoding the hex iv")
		}
	}
}
