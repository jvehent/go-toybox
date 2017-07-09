package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strings"
)

type ecdsaSignature struct {
	R, S *big.Int
}

func main() {
	sigBytes := decode("ZY03MVFlnWPo2fRH2s4FAjnzK92x93SGSWABSBR26lHG/7mShxaRum9Kfmi+E/o++HnNgSKK+L4sq9Bs2Plj/Q")
	if len(sigBytes)%2 != 0 {
		log.Fatal("invalid signature length, must be even, was", len(sigBytes))
	}
	// parse the signature
	fmt.Printf("\n# Trying to verify provided signature\n")
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(sigBytes[:len(sigBytes)/2])
	s.SetBytes(sigBytes[len(sigBytes)/2:])
	fmt.Printf("R: %s\nS: %s\n",
		r.String(),
		s.String())
}

func b64urlTob64(s string) string {
	// convert base64url characters back to regular base64 alphabet
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}
	return s
}

func decode(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(b64urlTob64(s))
	if err != nil {
		panic(err)
	}
	return b
}

func b64Tob64url(s string) string {
	// convert base64url characters back to regular base64 alphabet
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.TrimRight(s, "=")
	return s
}

func encode(b []byte) string {
	return b64Tob64url(base64.StdEncoding.EncodeToString(b))
}
