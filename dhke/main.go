package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"os"
)

func main() {

	primeGuy := "00c037c37588b4329887e61c2da3324b1ba4b81a63f9748fed2d8a410c2fc21b1232f0d3bfa024276cfd88448197aae486a63bfca7b8bf7754dfb327c7201f6fd17fd7fd74158bd31ce772c9f5f8ab584548a99a759b5a2c0532162b7b6218e8f142bce2c30d7784689a483e095e701618437913a8c39c3dd0d4ca3c500b885fe3"

	g := 2

	if len(os.Args) < 3 {
		fmt.Println("usage: ./main.exe [username] [password]")
		os.Exit(0)
	}

	var username string = os.Args[1]
	var password string = os.Args[2]
	var randomSeed uint64 = rand.Uint64()

	var privateKey string = username + ":" + password + ":" + fmt.Sprint(randomSeed)

	sum := sha256.New()
	sum.Write([]byte(privateKey))
	X := sum.Sum(nil)

	X = new(big.Int).Exp()

}
