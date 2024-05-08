package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	SERVER_HOST = "127.0.0.1"
	SERVER_PORT = "9001"
	SERVER_TYPE = "tcp"
)

func main() {

	args := os.Args

	// username:password combo
	// this should make it single simple string for dhke
	combo := args[1]
	mode := args[2]

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

		for {
			for i := 1023; i >= 0; i-- {
				if buffer[i] != 0 {
					buffer = buffer[:i+1]
					break
				}
			}
			dec_message := Decrypt(symkey, iv, buffer)

			fmt.Println(string(*dec_message))

			// message to AES encrypt
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("message to send: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			enc_message := Encrypt(symkey, iv, []byte(message))

			conn.Write(*enc_message)

			buffer = make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}
		}

	} else if mode == "c" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("ip to connect to: ")
		ip, _ := reader.ReadString('\n')
		ip = strings.TrimSpace(ip)

		conn = Setup_client(ip)

		//symkey, iv = Handshake(Y, conn, X)
		symkey, iv = Handshake(Y, conn, X)

		defer conn.Close()

		for {

			// message to AES encrypt
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("message to send: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			enc_message := Encrypt(symkey, iv, []byte(message))
			conn.Write(*enc_message)

			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}

			for i := 1023; i >= 0; i-- {
				if buffer[i] != 0 {
					buffer = buffer[:i+1]
					break
				}
			}

			dec_message := Decrypt(symkey, iv, buffer)

			fmt.Println(string(*dec_message))

		}
	}
}
