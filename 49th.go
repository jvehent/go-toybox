package main

import (
	"fmt"
	"math/big"
)

func main() {
	mersenne := new(big.Int)
	mersenne = mersenne.Exp(big.NewInt(2), big.NewInt(74207281), big.NewInt(0))
	mersenne = mersenne.Sub(mersenne, big.NewInt(1))
	fmt.Println(mersenne)
}
