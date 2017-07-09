package main

import "net/http"

func main() {
	_, err := http.Get("https://10000-sans.badssl.com")
	if err != nil {
		panic(err)
	}
}
