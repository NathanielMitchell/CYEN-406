package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	const message = "Team Name: Crypto Bros, Team Leader: Vito Mumphrey"

	msg := []byte(message)
	random := rand.Reader

	pemData, err := os.ReadFile("public.pem")
	checkErr(err)

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	cert, err := x509.ParsePKIXPublicKey(block.Bytes)
	checkErr(err)

	// message must be no longer than the length of the public modulus minus 11 bytes
	enc_message, err := rsa.EncryptPKCS1v15(random, cert.(*rsa.PublicKey), msg)
	checkErr(err)

	os.WriteFile("secret", enc_message, 777)
}
