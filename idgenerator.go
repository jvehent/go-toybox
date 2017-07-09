package main

import (
	"fmt"
	"sync"
	"time"
)

const maxiter int = 10000000

func main() {
	fmt.Printf("%.0f\n", GenID())
	for i := 0; i < maxiter; i++ {
		GenID()
		//fmt.Printf("%.0f\n", GenID())
	}
	fmt.Printf("%.0f\n", GenID())
}

type id struct {
	value float64
	sync.Mutex
}

var globalID id

func GenID() float64 {
	globalID.Lock()
	defer globalID.Unlock()
	if globalID.value < 1 {
		// if id hasn't been initialized yet, set it to number of seconds since
		// MIG's inception, plus one
		tmpid := int64(time.Since(time.Unix(1367258400, 0)).Seconds() + 1)
		tmpid = tmpid << 16
		globalID.value = float64(tmpid)
		return globalID.value
	} else {
		globalID.value++
		return globalID.value
	}
}
