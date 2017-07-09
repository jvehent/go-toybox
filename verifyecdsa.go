package main

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
)

type ecdsaSignature struct {
	R, S *big.Int
}

func main() {
	r := new(big.Int)
	s := new(big.Int)
	r.SetString("37549677831811818002499962379255599138575326196132998989397761117038543954426388403148891378515348145067971754600035", 10)
	s.SetString("527641238687683012158452496567809508917853977522148409893429482947055009339599346395932244045201871594264125012464", 10)

	keyBytes, err := base64.StdEncoding.DecodeString("MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE7oM/ewOhz6qtHyQhqJvT3SiefGPWqGwEUAZGVkuSIwvteVKrd8jnAjHYyCaYpIg9Vo10WnhXvm96L3KAbOE6Cyu3fMtKhZZIMf+Qqes9+66ae/NTeIWlDiGrjNeD+ClM")
	if err != nil {
		log.Fatal(err)
	}
	keyInterface, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	pubKey := keyInterface.(*ecdsa.PublicKey)

	md := sha512.New384()
	md.Write([]byte("Content-Signature:\x00cariboumaurice\n"))
	input := md.Sum(nil)

	fmt.Printf("signature verification: %t\n", ecdsa.Verify(pubKey, input, r, s))
}
