package main

import (
	"fmt"
	"os"
)

func main() {
	target := "lstat-link3"
	fi, err := os.Lstat(target)
	if err != nil {
		panic(err)
	}
	fmt.Println(fi.Name())
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		fmt.Println(os.Readlink(target))
	}

	fi, err = os.Stat(target)
	if err != nil {
		panic(err)
	}
	fmt.Println(fi.Name())
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		fmt.Println(os.Readlink(target))
	}
}
