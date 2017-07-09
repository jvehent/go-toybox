package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				rsa.GenerateKey(rand.Reader, 4096)
			}
		}()
	}
	go func() {
		i := 0
		for {
			i++
			fmt.Printf("%d ", i)
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(30 * time.Second)
	fmt.Println("setting GOMAXPROCS to 1")
	runtime.GOMAXPROCS(1)
	time.Sleep(30 * time.Second)
	fmt.Println("setting GOMAXPROCS to 10")
	runtime.GOMAXPROCS(10)
	time.Sleep(30 * time.Second)
	fmt.Println("setting GOMAXPROCS to 1")
	runtime.GOMAXPROCS(1)
}
