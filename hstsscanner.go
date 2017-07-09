package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, _ := http.Get("https://jve.linuxwall.info")
	if resp.Header.Get("Strict-Transport-Security") != "" {
		fmt.Println("HSTS Supported - Strict-Transport-Security", resp.Header.Get("Strict-Transport-Security"))
	}
}
