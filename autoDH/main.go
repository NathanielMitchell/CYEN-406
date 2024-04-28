package main

import (
    "os"

    "github.com/NathanMitchell/CYEN-406/dhke"
    "github.com/NathanMitchell/CYEN-406/simple_aes"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "9001"
	SERVER_TYPE = "tcp"
)

func main() {

	args := os.Args

	// flag for client/server usage
	// should be client|server
	mode := args[1]

	// username/password combo
	// should be in the format of username:password, for Nate :)
	// this should make it single simple string
	combo := args[2]

	// message to AES encrypt
	// idk how to format this
	message := []byte(args[3])

	// run the program based on the mode supplied
	switch mode {
	case "server":
        dhke
		server()
	case "client":
		client(message)
	}
}
