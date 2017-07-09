package main

import (
	"fmt"
	"time"
)

func main() {
	const somebytes = "\xab\xcd\xef\x12\x34\x56"
	fmt.Printf("%s\n", somebytes)
	time.Sleep(9999999 * time.Second)
}
