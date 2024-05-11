package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net"
)

func Setup_client(ip string) (conn net.Conn) {
	//establish connection
	conn, err := net.Dial(SERVER_TYPE, (ip + ":9001"))

	if err != nil {
		fmt.Println("error while trying to connect to the remote server")
		fmt.Println(err)
		return
	}
	return conn
}

func Send_data(message []byte, conn net.Conn) {

	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

}

func Handshake(Y string, conn net.Conn, X *big.Int) (symkey []byte, iv []byte) {

	symKeyChannel := make(chan []byte, 2)
	computationFinish := make(chan string)

	// send the pubkey to the other team

	go func(Y string, conn net.Conn, symKeyChannel <-chan []byte, computationFinish chan<- string) {
		_, err := conn.Write([]byte(Y))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
		}

		fmt.Println("Y has been sent")

		symKey := <-symKeyChannel
		iv := <-symKeyChannel

		fmt.Println("Sym key and IV have been received in go send")

		newSum := sha256.New()
		newSum.Write(append(symKey, iv...))
		hash := newSum.Sum(nil)

		_, err = conn.Write(hash)
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
		}

		fmt.Println("hash has been sent")

		computationFinish <- "Ready to return"

	}(Y, conn, symKeyChannel, computationFinish)

	go func(X *big.Int, conn net.Conn, symKeyChannel chan<- []byte) {
		buffer := make([]byte, 1024)
		bufferReset := make([]byte, 1024)

		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading: ", err.Error())
		}

		fmt.Println("other public key has been read")

		Yb := string(buffer)

		symKey, iv := DhkeGenerateSym(X, Yb)
		buffer = bufferReset

		symKeyChannel <- symKey
		symKeyChannel <- iv

		fmt.Println("sym key and iv have been sent in go rcv")

		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading: ", err.Error())
		}
		fmt.Printf("Server Response to Hash: %s\n", string(buffer))

		// resend the symkey and iv into the channel
		symKeyChannel <- symKey
		symKeyChannel <- iv
	}(X, conn, symKeyChannel)

	// Send_data([]byte(Y), conn)
	// fmt.Printf("Our Public Key: %s\n", Y)

	// // get their pubkey
	// buffer := make([]byte, 1024)
	// bufferReset := make([]byte, 1024)
	// _, err := conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("error reading:", err.Error())
	// }
	// fmt.Printf("this is what they Sent: %s\n", buffer)

	// symkey, iv = DhkeGenerateSym(X, string(buffer))

	// fmt.Printf("Sym Key: %s\n", hex.EncodeToString(symkey))
	// fmt.Printf("iv: %s\n", hex.EncodeToString(iv))
	// newSum := sha256.New()
	// // need to have it be utf-8 encoded for it to be compatible.
	// newSum.Write(append(symkey, iv...))
	// hash := newSum.Sum(nil)

	// fmt.Printf("Hash of Key+IV: %s\n", hex.EncodeToString(hash))

	// // send the challenge hash
	// _, err = conn.Write(hash)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// buffer = bufferReset
	// _, err = conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("Server Response to Hash: %s\n", string(buffer))

	fmt.Println(<-computationFinish)

	symKey := <-symKeyChannel
	IV := <-symKeyChannel

	return symKey, IV
}
