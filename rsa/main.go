package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	privateKey,err:= rsa.GenerateKey(rand.Reader,2048)
	if err != nil {
		log.Fatal(err)
	}

	// the public key is a part of *rsa.PrivateKey struct
	publickey := privateKey.PublicKey

	modulesBytes := base64.StdEncoding.EncodeToString(privateKey.N.Bytes())
	privateExponentBytes := base64.StdEncoding.EncodeToString(privateKey.D.Bytes())
	fmt.Println(modulesBytes)
	fmt.Println(privateExponentBytes)
	fmt.Println(publickey.E)

	encrytedBytes ,err :=rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publickey,
		[]byte("super secret message"),
		nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("encrypted bytes: ",encrytedBytes)

	// the first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// the OEAPOptions in the end signify that we encrypted the data using OEAP,and that we used
	// SHA256 to hash the input

	decryptedBytes , err := privateKey.Decrypt(nil,encrytedBytes,&rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("decrypted message: ",string(decryptedBytes))

	msg := []byte("verifiable message")

	// before signing ,we need to hash our message
	// the hash is what we actually sign
	msgHash := sha256.New()
	_,err =msgHash.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	msgHashSUm := msgHash.Sum(nil)

	signature , err := rsa.SignPSS(rand.Reader,privateKey,crypto.SHA256,msgHashSUm,nil)
	if err != nil {
		log.Fatal(err)
	}

	err = rsa.VerifyPSS(&publickey,crypto.SHA256,msgHashSUm,signature,nil)
	if err != nil {
		fmt.Println("coud not verify signature: ",err)
		return
	}

	fmt.Println("signaure verified")
}
