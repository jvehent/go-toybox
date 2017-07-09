package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	switch os.Args[1] {
	case "generate":
		salt := getSalt()
		dk := pbkdf2.Key(
			[]byte(os.Args[2]), // read password from stdin
			salt,               // add 16 bytes of salt
			65536,              // perform 4096 iterations
			32,                 // desired output length
			sha256.New,         // hashing algorithm
		)
		fmt.Printf("hash=%X\nsalt=%X\n", dk, salt)
	case "verify":
		salt, err := hex.DecodeString(os.Args[4])
		if err != nil {
			log.Fatal("failed to hex decode hash: %v", err)
		}
		dk := pbkdf2.Key(
			[]byte(os.Args[2]), // read password from stdin
			salt,               // add 16 bytes of salt
			65536,              // perform 4096 iterations
			32,                 // desired output length
			sha256.New,         // hashing algorithm
		)
		if fmt.Sprintf("%X", dk) == os.Args[3] {
			fmt.Println("ok")
		}
	default:
		fmt.Printf(`pbkdf2 generator & verifier.
usage:
./pbkdf2 generate somepassword
./pbkdf2 verify somepassword <hex_encoded_hash> <hex_encoded_salt>
`)
	}
}

func getSalt() []byte {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal("error getting salt:", err)
	}
	return salt
}
