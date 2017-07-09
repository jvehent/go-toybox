package main

import (
	"code.google.com/p/go.crypto/openpgp"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func main() {
	secringFile, err := os.Open("/tmp/testgolangopenpgp/secring.gpg")
	if err != nil {
		panic(err)
	}
	defer secringFile.Close()
	keyring, err := openpgp.ReadKeyRing(secringFile)
	if err != nil {
		err = fmt.Errorf("Keyring access failed: '%v'", err)
		panic(err)
	}
	fmt.Printf("found %d entities in keyring\n", len(keyring))
	for _, entity := range keyring {
		fingerprint := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
		fmt.Println("reading entity with fingerprint", fingerprint)
	}
}
