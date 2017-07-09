package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

func main() {
	dk, err := scrypt.Key([]byte("some password"), []byte("some salt"), 16384, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%X\n", dk)
}
