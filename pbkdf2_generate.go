package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func main() {
	salt := getSalt()
	dk := pbkdf2.Key(
		[]byte(os.Args[1]),
		salt,
		65536,
		32,
		sha256.New,
	)
	fmt.Printf("hash=%X\nsalt=%X\n", dk, salt)
}

func getSalt() []byte {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal("error getting salt:", err)
	}
	return salt
}
