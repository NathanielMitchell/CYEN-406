package main

import (
	"fmt"
	"net"
)

func client(message []byte) {
	//establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
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
