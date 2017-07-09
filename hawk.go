package main

import (
	"fmt"

	"github.com/mozilla-services/hawk-go"
)

func main() {
	auth, err := hawk.ParseRequestHeader(`Hawk id="dh37fgj492je", ts="1353832234", nonce="j4h3g2", ext="some-app-ext-data", mac="6R4rV5iE+NPoym+WwjeHzjAGXUtLNIxmo1vpMofpLAE="`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *auth)
}
