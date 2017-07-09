package main

import (
	"crypto/rand"
	"fmt"
)

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func main() {
	fmt.Println(randToken())
}
