package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	addrs, err := net.LookupHost(name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("addr 0 => '%s'\n", addrs[0])
	fmt.Printf("addr array => '%s'\n", addrs)
}
