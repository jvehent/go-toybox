package main

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func main() {
	PrivateKey := "****************************"
	rawkey, err := base64.StdEncoding.DecodeString(PrivateKey)
	if err != nil {
		panic(err)
	}
	ecdsaPrivKey, err := x509.ParseECPrivateKey(rawkey)
	if err != nil {
		panic(err)
	}
	pubkeybytes, err := x509.MarshalPKIXPublicKey(ecdsaPrivKey.Public())
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(pubkeybytes))
}
