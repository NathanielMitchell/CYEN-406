package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
    "fmt"
	"log"
	"os"
)

var message = "Team Name: Crypto Bros, Team Leader: Vito Mumphrey"

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Takes an argument of string for the public key filename
// encrypts the message
func encrypt(filename string) {

	msg := []byte(message)
	random := rand.Reader

	pemData, err := os.ReadFile("public.pem")
	checkErr(err)

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block")
	}

	cert, err := x509.ParsePKIXPublicKey(block.Bytes)
	checkErr(err)

	// message must not be > length of the (pub modulus-11) bytes
	enc_message, err := rsa.EncryptPKCS1v15(random, cert.(*rsa.PublicKey), msg)
	checkErr(err)

	os.WriteFile("secret", enc_message, 0666)

    fmt.Println("Encrypted Message...")
    fmt.Println(enc_message)

    return
}

// Takes in the private keys filename
// decrypts the secret message which is in "secret"
func decrypt(filname string) {

    secret_msg, err := os.ReadFile("secret")
    checkErr(err)

	random := rand.Reader
    
    pemData, err := os.ReadFile("private.pem")
	checkErr(err)

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode PEM block")
	}

	cert, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	checkErr(err)

    dec_message, err := rsa.DecryptPKCS1v15(random, cert.(*rsa.PrivateKey), secret_msg)
    checkErr(err)

    fmt.Println("Decrypted Message...")
    fmt.Println(string(dec_message))

    return
}

func main () {
    args := os.Args
    
    if len(args) != 3 {
        fmt.Println("./main mode keyfilename")
        return
    }

    mode := args[1]

    switch mode {
    case "encrypt":
        encrypt(args[2]) 
    case "decrypt":
        decrypt(args[2])
    default:
        fmt.Println("./main mode keyfilename")
        return
    }

    return
}
