package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net"
)

func Setup_client(ip string) (conn net.Conn) {
	//establish connection
	conn, err := net.Dial(SERVER_TYPE, ip)

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
	// send the pubkey to the other team
	Send_data([]byte(Y), conn)

	// get their pubkey
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("error reading:", err.Error())
	}

	symkey, iv = DhkeGenerateSym(X, string(buffer))

	newSum := sha256.New()
	// need to have it be utf-8 encoded for it to be compatible.
	newSum.Write(append(symkey, iv...))
	hash := newSum.Sum(nil)

	// send the challenge hash
	_, err = conn.Write(hash)
	if err != nil {
		fmt.Println(err)
	}

	return symkey, iv
}