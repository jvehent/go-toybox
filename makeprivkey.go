package main

import (
	"bytes"
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/armor"
	"fmt"
)

func main() {
	ent, err := openpgp.NewEntity("bob", "Bob's key", "bob@example.net", nil)
	if err != nil {
		panic(err)
	}
	pkbuf := bytes.NewBuffer(nil)
	err = ent.SerializePrivate(pkbuf, nil)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	ewrbuf, err := armor.Encode(buf, "PGP PRIVATE KEY BLOCK", nil)
	if err != nil {
		panic(err)
	}
	_, err = ewrbuf.Write(pkbuf.Bytes())
	if err != nil {
		panic(err)
	}
	ewrbuf.Close()
	fmt.Printf("%s\n", buf.Bytes())

	// serialize the public key
	pkbuf = bytes.NewBuffer(nil)
	err = ent.Serialize(pkbuf)
	if err != nil {
		panic(err)
	}
	buf = bytes.NewBuffer(nil)
	ewrbuf, err = armor.Encode(buf, openpgp.PublicKeyType, nil)
	if err != nil {
		panic(err)
	}
	_, err = ewrbuf.Write(pkbuf.Bytes())
	if err != nil {
		panic(err)
	}
	ewrbuf.Close()
	fmt.Printf("%s\n", buf.Bytes())

}
