package autodh

import (
	"fmt"
	"net"
	"os"

    "github.com/NathanielMitchell/CYEN-406/dhke"
)

func server(pubkey string) {
    fmt.Printf("Server Running on SERVER_HOST:SERVER_PORT...\n", SERVER_HOST, SERVER_PORT)

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection, pubkey)
	}
}

// This function is used to start the diffy-helman progress
// returns the public key of the server
func processClient(connection net.Conn, pubkey string) (client_connection net.Conn) {
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	_, err = connection.Write([]byte(pubkey))
	
    return client_connection
}
