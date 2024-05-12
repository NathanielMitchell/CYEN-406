package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
)

func Server(X *big.Int, Y string) (symkey []byte, iv []byte, connection net.Conn) {
	fmt.Printf("Server Running on %s:%s...\n", SERVER_HOST, SERVER_PORT)

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

	symkey, iv = ProcessClient(connection, X, Y)

	return symkey, iv, connection
}

// This function is used to start the diffy-helman progress
// returns the public key of the server
func ProcessClient(connection net.Conn, X *big.Int, Y string) (symkey []byte, iv []byte) {

	symKeyChannel := make(chan []byte, 2)
	finishChannel := make(chan string)

	go func(symKeyChannel chan<- []byte, X *big.Int, connection net.Conn, finishChannel chan string) {
		buffer := make([]byte, 1024)
		bufferReset := make([]byte, 1024)

		_, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error Reading: ", err.Error())
		}

		yb := string(buffer)
		buffer = bufferReset

		symKey, iv := DhkeGenerateSym(X, yb)

		_, err = connection.Read(buffer)
		if err != nil {
			fmt.Println("Error Reading: ", err.Error())
		}

		otherTeamHash := buffer[:32]

		newSum := sha256.New()
		// need to have it be utf-8 encoded for it to be compatible.
		newSum.Write(append(symKey, iv...))
		hash := newSum.Sum(nil)

		symKeyChannel <- hash
		symKeyChannel <- otherTeamHash

		symKeyChannel <- symKey
		symKeyChannel <- iv
	}(symKeyChannel, X, connection, finishChannel)

	go func(Y string, connection net.Conn, symKeyChannel <-chan []byte, finishChannel chan<- string) {
		_, err := connection.Write([]byte(Y))
		if err != nil {
			fmt.Println("Error Writing: ", err.Error())
		}

		hash := <-symKeyChannel
		otherTeamHash := <-symKeyChannel

		if strings.Compare(string(hash), string(otherTeamHash)) == 0 {
			connection.Write([]byte("you good blud"))
		} else {
			connection.Write([]byte("blud...we're cooked"))
		}

		finishChannel <- "ready"

	}(Y, connection, symKeyChannel, finishChannel)
	// buffer := make([]byte, 1024)
	// bufferReset := make([]byte, 1024)
	// _, err := connection.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// }

	// // the other teams public key
	// symkey, iv = DhkeGenerateSym(X, string(buffer))

	// newSum := sha256.New()
	// // need to have it be utf-8 encoded for it to be compatible.
	// newSum.Write(append(symkey, iv...))
	// hash := newSum.Sum(nil)

	// // give the other team our pubkey
	// _, err = connection.Write([]byte(Y))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// buffer = bufferReset
	// _, err = connection.Read(buffer)
	// if err != nil {
	// 	fmt.Println("error reading:", err.Error())
	// }

	// if strings.Compare(string(hash), string(buffer[:32])) == 0 {
	// 	connection.Write([]byte("you good blud"))
	// } else {
	// 	connection.Write([]byte("blud...we're cooked"))
	// }

	<-finishChannel
	symKey := <-symKeyChannel
	IV := <-symKeyChannel

	return symKey, IV
}
