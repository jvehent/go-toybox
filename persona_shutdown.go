package main

import (
	"fmt"
	"time"
)

func main() {
	t, _ := time.Parse("2006-01-02", "2016-12-01")
	fmt.Println("Persona shuts down on", t.String(), "in", t.Sub(time.Now()).Seconds()/86400, "days")
}
