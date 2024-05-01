package main

import (
	"fmt"
	"net"
)

func client(message []byte, ip string) {
	//establish connection
	connection, err := net.Dial(SERVER_TYPE, ip+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("error while trying to connect to the remote server")
        return
	}

	// send some data
	_, err = connection.Write(message)
	buffer := make([]byte, len(message))
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Println("Received: ", string(buffer[:mLen]))
	defer connection.Close()
}
