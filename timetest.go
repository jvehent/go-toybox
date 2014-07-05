package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now().UTC().Format(time.RFC3339Nano)
	fmt.Println(t)
}
