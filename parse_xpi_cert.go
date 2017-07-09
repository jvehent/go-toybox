package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fullsailor/pkcs7"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	p7, err := pkcs7.Parse(data)
	if err != nil {
		panic(err)
	}
	for _, cert := range p7.Certificates {
		fmt.Println(cert.Subject.CommonName)
	}
	for _, siginfo := range p7.Signers {
		fmt.Printf("%X\n", siginfo.EncryptedDigest)
	}
}
