package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func main() {
	for _, size := range []int{1024, 2048, 3072, 4096} {
		priv, _ := rsa.GenerateKey(rand.Reader, size)
		fmt.Printf("\n%d\nN: %s\nE: %d\nD: %s\np: %s\nq: %s\n\n",
			size,
			priv.PublicKey.N.String(),
			priv.PublicKey.E,
			priv.D.String(),
			priv.Primes[0].String(),
			priv.Primes[1].String())
	}
}
