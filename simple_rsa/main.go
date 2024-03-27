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

var DEBUG bool = false

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
// Takes an argument of string for the public key filename
// encrypts the message
// prints to stdout and file named 'secret'
func encrypt(filename string, msg_filename string) {

    msg, err := os.ReadFile(msg_filename)
    checkErr(err)

	random := rand.Reader

	pemData, err := os.ReadFile(filename)
	checkErr(err)

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block")
	}

	cert, err := x509.ParsePKIXPublicKey(block.Bytes)
	checkErr(err)
  f, err := os.OpenFile("secret", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666);

  key_size := (cert.(*rsa.PublicKey).N.BitLen() + 7) / 8
  chunk_size := key_size - 11
  fmt.Printf("Key Size: %d, Chunk size: %d\n", key_size, chunk_size)
  for i := 0;  i < len(msg); i += chunk_size{  
    if DEBUG {fmt.Printf("i: %d\n", i)}
    if len(msg) < (i + chunk_size){
      enc_message, err := rsa.EncryptPKCS1v15(random, cert.(*rsa.PublicKey), msg[i:len(msg)])
	    checkErr(err)
	    if DEBUG {fmt.Println("Writing to ./secret...")}
      f.Write(enc_message);
	    if DEBUG {
        fmt.Println("Encrypted Message...")
	      fmt.Println(enc_message)
      }
    } else{
      enc_message, err := rsa.EncryptPKCS1v15(random, cert.(*rsa.PublicKey), msg[i:(i + chunk_size)])
	    checkErr(err)
	    if DEBUG {fmt.Println("Writing to ./secret...")}
      f.Write(enc_message)
      if DEBUG {
	      fmt.Println("Encrypted Message...")
	      fmt.Println(enc_message)
      }
    }
    // figure out scope to put this outside the logic
    /*
	    fmt.Println("Writing to ./secret...")
	    os.WriteFile("secret", enc_message, 0666)
	    fmt.Println("Encrypted Message...")
	    fmt.Println(enc_message)
    */
    }
}

// Takes in the private keys filename
// decrypts the secret message which is in "secret"
func decrypt(filname string, secret_msg_filname string) {

    secret_msg, err := os.ReadFile(secret_msg_filname)
    checkErr(err)

	pemData, err := os.ReadFile(filname)
	checkErr(err)

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode PEM block")
	}

	cert, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	checkErr(err)

  key_size := (cert.(*rsa.PrivateKey).PublicKey.N.BitLen() + 7) / 8
  chunk_size := key_size
  if DEBUG {fmt.Printf("Key Size: %d, Chunk size: %d\n", key_size, chunk_size)}
  for i := 0;  i < len(secret_msg); i += chunk_size{
    if DEBUG {fmt.Printf("i: %d\n", i)}
    if len(secret_msg) < (i + chunk_size){
      dec_message, err := rsa.DecryptPKCS1v15(nil, cert.(*rsa.PrivateKey), secret_msg[i:len(secret_msg)])
	    checkErr(err)
      if DEBUG{
	    fmt.Println("Decrypted Message...")
      }
	    fmt.Println(string(dec_message))
    } else{
      dec_message, err := rsa.DecryptPKCS1v15(nil, cert.(*rsa.PrivateKey), secret_msg[i:(i + chunk_size)])
	    checkErr(err)
      if DEBUG {
	      fmt.Println("Decrypted Message...")
      }
	    fmt.Println(string(dec_message))
    }
    /*
	  fmt.Println("Decrypted Message...")
	  fmt.Println(string(dec_message))
    */
  }
}

func main() {
	args := os.Args

	if len(args) != 4 {
		fmt.Println("./main mode[encrypt|decrypt] './keyfilename' './msg_file,txt'")
		return
	}

	mode := args[1]

	switch mode {
	case "encrypt":
		encrypt(args[2], args[3])
	case "decrypt":
		decrypt(args[2], args[3])
	default:
		fmt.Println("./main mode[encrypt|decrypt] './keyfilename' './msg_file,txt'")
	}

}
