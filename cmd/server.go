package main

import (
	"fmt"
	"net"
	"os"

	"github.com/NathanielMitchell/CYEN-406/dhke"
)

func server(combo string) (con chan net.Conn) {
	fmt.Printf("Server Running on SERVER_HOST:SERVER_PORT...\n", SERVER_HOST, SERVER_PORT)

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)

    connections := make(chan net.Conn)
	// listen for connections in go routine
	go func() {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("client connected %s\n", connection.RemoteAddr())
        go processClient(connection, combo, connections)
	}()

    return connections
}

// This function is used to start the diffy-helman progress
// returns the public key of the server
func processClient(connection net.Conn, combo string, connections <-chan net.Conn)  {
	buffer := make([]byte, 1024)
	otherPublicKey, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

    // the other teams public key
    resultGuy := dhke.Dhke(combo, otherPublicKey)

	// give the other team our pubkey
	_, err = connection.Write([]byte(resultGuy.ourPublicKey))

}
