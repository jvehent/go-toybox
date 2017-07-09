package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	premier_nombre, _ := strconv.Atoi(os.Args[1])
	deuxieme_nombre, _ := strconv.Atoi(os.Args[2])
	somme := premier_nombre + deuxieme_nombre
	fmt.Printf("la somme de %d + %d est %d\n", premier_nombre, deuxieme_nombre, somme)
}
