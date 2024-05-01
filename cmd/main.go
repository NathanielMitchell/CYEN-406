package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "9001"
	SERVER_TYPE = "tcp"
)

func main() {

	args := os.Args

	// username/password combo
	// this should make it single simple string for dhke
	combo := args[1]
	mode := args[2]
	ip := args[3]

	X, Y := DhkeGeneratePubKey(combo)

	var symkey []byte
	var iv []byte
	var conn net.Conn

	if mode == "s" {
		symkey, iv, conn = Server(X, Y)

		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		for true {
			dec_message := Decrypt(buffer, symkey, iv)

			fmt.Println(dec_message)

			// message to AES encrypt
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("message to send: ")
			message, _ := reader.ReadString('\n')

			enc_message := Encrypt([]byte(message), symkey, iv)

			conn.Write([]byte(*enc_message))

			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}
		}

	} else if mode == "c" {
		conn := Setup_client(ip)
		symkey, iv = Handshake(Y, conn, X)

		for true {
			// message to AES encrypt
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("message to send: ")
			message, _ := reader.ReadString('\n')

			enc_message := Encrypt([]byte(message), symkey, iv)
			conn.Write([]byte(*enc_message))

			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}

			dec_message := Decrypt(buffer, symkey, iv)

			fmt.Println(dec_message)
		}
	}
}
