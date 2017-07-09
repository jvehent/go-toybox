package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := []byte("ThisIsSecret")
	t1 := time.Now()
	iter := 10000
	for i := 0; i < iter; i++ {
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		// Comparing the password with the hash
		err = bcrypt.CompareHashAndPassword(hashedPassword, password)
		if err != nil {
			panic("password comparison failed")
		}
	}
	duration := time.Now().Sub(t1)
	fmt.Printf("total time: %s\n", duration.String())
	singlepass := duration.Nanoseconds() / int64(iter*2)
	fmt.Printf("single roundtrip: %d ms\n", singlepass/1000000)
}
