package main

import (
	"net"
)

// struct used to hold data for a team
type connectionHandler struct {
	symKey     []byte
	iV         []byte
	connection net.Conn
}
