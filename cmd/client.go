package main

import (
	"fmt"
	"net"
)

func setup_client(message []byte, ip string) (conn net.Conn) {
	//establish connection
	connection, err := net.Dial(SERVER_TYPE, ip+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("error while trying to connect to the remote server")
		return
	}
	return conn
}

func send_data(message []byte, conn net.Conn) () {
	// send some data
	_, err = connection.Write(message)

	buffer := make([]byte, len(message))
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

}

// read data from the channel buffer and print to stdout
func recieve_data(buffer chan, mlen) {
	fmt.Println("Received: ", string(buffer[:mlen]))
}
