package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {
	// generate an ecdsa private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// sign random data with it
	hashed := []byte("testing some string signing")
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)
	if err != nil {
		panic(err)
	}

	// extract the pubkey into an hexadecimal string
	pubkey := priv.Public()
	pubkeyhex := hex.EncodeToString(
		elliptic.Marshal(pubkey.(elliptic.Curve),
			pubkey.(*ecdsa.PublicKey).X,
			pubkey.(*ecdsa.PublicKey).Y),
	)
	fmt.Println("got pubkey", pubkeyhex)

	// reload the pubkey from an hexadecimal string
	pubkeybytes, err := hex.DecodeString(pubkeyhex)
	if err != nil {
		panic(err)
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pubkeybytes)
	var pubkey2 ecdsa.PublicKey
	pubkey2.Curve = elliptic.P256()
	pubkey2.X = x
	pubkey2.Y = y

	// verify the signature
	fmt.Printf("Signature is %t\n", ecdsa.Verify(&pubkey2, hashed, r, s))
}
