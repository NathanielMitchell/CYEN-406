package main

import (
	"os"

	"../dhke"
	"../simple_aes"
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
		dhke.dhke(combo)
		server()
	case "client":
		message := simple_aes.simple_aes(message)
		client(message)
	}
}
