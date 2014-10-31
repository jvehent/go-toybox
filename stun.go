package main

import (
	"fmt"

	"github.com/ccding/go-stun/stun"
)

func main() {
	stun.SetServerHost("v.mozilla.com", 3478)
	nat, host, err := stun.Discover()
	if err != nil {
		panic(err)
	}
	fmt.Println(host.Ip(), nat)
}
