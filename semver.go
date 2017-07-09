package main

import (
	"fmt"
	"log"

	"github.com/blang/semver"
)

func main() {
	v1, err := semver.Make("1.0.0")
	if err != nil {
		log.Fatal(err)
	}
	v2, err := semver.Make("2.0.1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%t\n", v1.Compare(v2))
}
