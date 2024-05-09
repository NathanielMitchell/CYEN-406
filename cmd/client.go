package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	// send the pubkey to the other team
	Send_data([]byte(Y), conn)
	fmt.Printf("Our Public Key: %s\n", Y)

	// get their pubkey
	buffer := make([]byte, 1024)
	bufferReset := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("error reading:", err.Error())
	}
	fmt.Printf("this is what they Sent: %s\n", buffer)

	symkey, iv = DhkeGenerateSym(X, string(buffer))

	fmt.Printf("Sym Key: %s\n", hex.EncodeToString(symkey))
	fmt.Printf("iv: %s\n", hex.EncodeToString(iv))
	newSum := sha256.New()
	// need to have it be utf-8 encoded for it to be compatible.
	newSum.Write(append(symkey, iv...))
	hash := newSum.Sum(nil)

	fmt.Printf("Hash of Key+IV: %s\n", hex.EncodeToString(hash))

	// send the challenge hash
	_, err = conn.Write(hash)
	if err != nil {
		fmt.Println(err)
	}
	buffer = bufferReset
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Server Response to Hash: %s\n", string(buffer))

	return symkey, iv
}
