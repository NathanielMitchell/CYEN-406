package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
)

// function takes in a string of "username:password"
func DhkeGeneratePubKey(combo string) (X *big.Int, Y string) {
	primeGuy, _ := new(big.Int).SetString("00c037c37588b4329887e61c2da3324b1ba4b81a63f9748fed2d8a410c2fc21b1232f0d3bfa024276cfd88448197aae486a63bfca7b8bf7754dfb327c7201f6fd17fd7fd74158bd31ce772c9f5f8ab584548a99a759b5a2c0532162b7b6218e8f142bce2c30d7784689a483e095e701618437913a8c39c3dd0d4ca3c500b885fe3", 16)
	g := big.NewInt(2)

	randomSeed := rand.Uint64()

	privateKey := combo + ":" + fmt.Sprint(randomSeed)

	sum := sha256.New()
	sum.Write([]byte(privateKey))
	// needs to be a []byte array
	// this generates the hash as a []byte
	h := sum.Sum(nil)
	X = new(big.Int).SetBytes(h)

	Y = new(big.Int).Exp(g, X, primeGuy).Text(16)

	return X, Y
}

func DhkeGenerateSym(X *big.Int, otherTeamPublicKey string) (symkey []byte, iv []byte) {
	primeGuy, _ := new(big.Int).SetString("00c037c37588b4329887e61c2da3324b1ba4b81a63f9748fed2d8a410c2fc21b1232f0d3bfa024276cfd88448197aae486a63bfca7b8bf7754dfb327c7201f6fd17fd7fd74158bd31ce772c9f5f8ab584548a99a759b5a2c0532162b7b6218e8f142bce2c30d7784689a483e095e701618437913a8c39c3dd0d4ca3c500b885fe3", 16)

	otherPublicKey := new(big.Int)
	_, _ = otherPublicKey.SetString(otherTeamPublicKey, 16)

	K := new(big.Int).Exp(otherPublicKey, X, primeGuy)

	newSum := sha256.New()
	// need to have it be utf-8 encoded for it to be compatible.
	newSum.Write([]byte(K.Text(10)))
	outputSymmetricKey := newSum.Sum(nil)

	newNewSum := md5.New()
	newSum.Write([]byte(K.Text(10)))
	outputIV := newNewSum.Sum(nil)

	return outputSymmetricKey, outputIV
}
