package main

import (
	"fmt"
	"os"
)

func main() {
	fd, _ := os.Open("/usr/bin/yelp")
	fmt.Println(fd.Name())
	fd.Close()
}
