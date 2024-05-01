package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
)

func server(X *big.Int, Y string) (symkey []byte, iv []byte, connection net.Conn) {
	fmt.Printf("Server Running on SERVER_HOST:SERVER_PORT...\n", SERVER_HOST, SERVER_PORT)

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)

	// listen for connections in go routine
    connection, err = server.Accept()
    if err != nil {
        fmt.Println("Error accepting: ", err.Error())
        os.Exit(1)
    }
    fmt.Printf("client connected %s\n", connection.RemoteAddr())
    
    symkey, iv = processClient(connection, X, Y)

    return symkey, iv, connection
}

// This function is used to start the diffy-helman progress
// returns the public key of the server
func processClient(connection net.Conn, X *big.Int, Y string) (symkey []byte, iv []byte)  {
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

    // the other teams public key
    symkey, iv = DhkeGenerateSym(X, string(buffer))

	newSum := sha256.New()
	// need to have it be utf-8 encoded for it to be compatible.
	newSum.Write(append(symkey, iv...))
	hash := newSum.Sum(nil)

	// give the other team our pubkey
	_, err = connection.Write([]byte(Y))
    
    _, err = connection.Read(buffer)
	if err != nil {
		fmt.Println("error reading:", err.Error())
	}

    if strings.Compare(string(hash), string(buffer)) == 0 {
        connection.Write([]byte("you good blud"))
    } else {
        connection.Write([]byte("blud...we're cooked"))
    }

    return symkey, iv
}
