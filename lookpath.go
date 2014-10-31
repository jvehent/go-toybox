package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	path, err := exec.LookPath("lsb_release")
	if err != nil {
		log.Fatal("lsb_release is missing")
	}
	fmt.Printf("lsb_release is available at %s\n", path)
}
