package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	target, err := filepath.EvalSymlinks("/sbin/init")
	if err != nil {
		panic(err)
	}
	fmt.Println(target)
}
