package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.NumCPU(), "cpus; ", runtime.GOMAXPROCS(-1), "gomaxprocs")
}
