package main

import (
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/armor"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func main() {
	sigReader := strings.NewReader(sig)
	sigBlock, err := armor.Decode(sigReader)
	if err != nil {
		panic(err)
	}
	if sigBlock.Type != openpgp.SignatureType {
		panic("not a signature type")
	}

	dataReader := strings.NewReader(data)

	krfd, err := os.Open("/home/ulfr/.gnupg/pubring.gpg")
	if err != nil {
		panic(err)
	}
	defer krfd.Close()

	keyring, err := openpgp.ReadKeyRing(krfd)
	if err != nil {
		panic(err)
	}
	entity, err := openpgp.CheckDetachedSignature(
		keyring, dataReader, sigBlock.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("valid signature from key %s\n",
		hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
}

var sig string = `-----BEGIN PGP SIGNATURE-----

iQEcBAABCAAGBQJUWOB4AAoJEKPWUhc7dj6PvMMH/im1r3YaTfc+l0JjcPIkXIKI
FQApAN0RqbLG5845oiJ9uDPn/1wQP8qzxXldOUXOpHChy7EXf8v8TuwHlNeJAX5r
b8zooX6gPtSrDd25okOda9hdBK7h4Niy0V2Lwx+DCikO0xmTA6Ftxg7n3Z+chfey
//OPpVjhcorythf2aylF3XK3ZzJcOg5aX+N6OubYqd1nztIo313NS4Q//ctvhhCN
89dNhuLHVrb3NTLFR/WXXBleyhBF1omNmTIR+F4V9bDdH1gsvFyMQIPOIJMeagvb
9Jmr/nQyVGZVeqDkMfh6iw4Tjs/soCuIN4BmVH13PmtdtZuDVwy5sz3Qqc3h+gQ=
=1r0q
-----END PGP SIGNATURE-----`

var data string = "2014-11-04T12:39:20.0Z;1825922807490630059\n"
