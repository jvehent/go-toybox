package main

import (
	"fmt"
	"os"
)

func main() {
	char := fmt.Sprintf("%c", os.Args[0][0])
	fmt.Println(char)
}
