package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strings"
)

var PrivateKey string = "MIGkAgEBBDAzX2TrGOr0WE92AbAl+"
var EncryptionKey string = "MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE4k3FmG7dFoOt3Tuzl76abTRtK8sb/"

//var ContentSignature string = "XBKzej3i6TAFZc3VZsuCekn-4dYWJBE4-b3OOtKrOV-JIzIvAnAhnOV1aj-kEm07kh-FciIxV-Xk2QUQlRQzHO7oW7E4mXkMKkbbAcvL0CFrItTObhfhKnBnpAE9ql1O"
var ContentSignature string = "VhKPP_rOtJUQVhpWMBrG3w5B6CMKU8XoksL0pyJg39jj90NFJLXP07gpL_FiaAcfUK9csW4KKxp83sDYyrcnT_n4HuoBBps_DJs-hN4BYzs40kY_F34Vp0q5ph9ofobw"
var Content string = `<!DOCTYPE HTML>
<html>
<!-- https://bugzilla.mozilla.org/show_bug.cgi?id=1226928 -->
<head>
  <meta charset="utf-8">
  <title>Testpage for bug 1226928</title>
</head>
<body>
  Just a fully good testpage for Bug 1226928<br/>
</body>
</html>`

type ecdsaSignature struct {
	R, S *big.Int
}

func main() {
	// parse the public key
	keyInterface, err := x509.ParsePKIXPublicKey(decode(EncryptionKey))
	if err != nil {
		panic(err)
	}
	pubKey := keyInterface.(*ecdsa.PublicKey)
	fmt.Printf("Curve: %s (%d bits)\n\tb64: %s\n\tP: %s\n\tN: %s\n\tB: %s\n\tGx: %s\n\tGy: %s\n\tX: %s\n\tY: %s\n",
		pubKey.Params().Name, pubKey.Params().BitSize, b64urlTob64(EncryptionKey),
		pubKey.Params().P.String(), pubKey.Params().N.String(), pubKey.Params().B.String(),
		pubKey.Params().Gx.String(), pubKey.Params().Gy.String(),
		pubKey.X.String(), pubKey.Y.String())

	// calculate the hash
	md := sha512.New384()
	md.Write([]byte("Content-Signature:" + "\x00" + Content))
	hash := md.Sum(nil)
	fmt.Printf("Content:\n\tsha384: %X\n----- begin content -----\n%s\n----- end content -----\n", hash, Content)

	sigBytes := decode(ContentSignature)
	if len(sigBytes)%2 != 0 {
		log.Fatal("invalid signature length, must be even, was", len(sigBytes))
	}
	// parse the signature
	fmt.Printf("\n# Trying to verify provided signature\n")
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(sigBytes[:len(sigBytes)/2])
	if err != nil {
		panic(err)
	}
	s.SetBytes(sigBytes[len(sigBytes)/2:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Signature:\n\tb64: %s\n\tR: %s\n\tS: %s\n",
		b64urlTob64(ContentSignature),
		r.String(),
		s.String())

	// verify the content signature
	fmt.Printf("Signature is %t\n", ecdsa.Verify(pubKey, hash, r, s))

	fmt.Printf("\n# Trying to sign & verify ourselves\n")
	// parse the private key
	privKey, err := x509.ParseECPrivateKey(decode(PrivateKey))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Private Key:\n\tD: %s\n", privKey.D.String())
	ecdsaSig := new(ecdsaSignature)
	ecdsaSig.R, ecdsaSig.S, err = ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Signature is %t\n", ecdsa.Verify(pubKey, hash, ecdsaSig.R, ecdsaSig.S))
	sigAsn1Bytes, err := asn1.Marshal(ecdsaSignature{ecdsaSig.R, ecdsaSig.S})
	if err != nil {
		panic(err)
	}
	sigWebCryptoBytes := make([]byte, len(ecdsaSig.R.Bytes())+len(ecdsaSig.S.Bytes()))
	copy(sigWebCryptoBytes[:len(ecdsaSig.R.Bytes())], ecdsaSig.R.Bytes())
	copy(sigWebCryptoBytes[len(ecdsaSig.R.Bytes()):], ecdsaSig.S.Bytes())
	fmt.Printf("Signature:\n\tASN1 b64: %s\n\tASN1 b64url: %s\n\twebcrypto b64: %s\n\twebcrypto b64url: %s\n\tR: %s\n\tS: %s\n",
		base64.StdEncoding.EncodeToString(sigAsn1Bytes), encode(sigAsn1Bytes),
		base64.StdEncoding.EncodeToString(sigWebCryptoBytes), encode(sigWebCryptoBytes),
		ecdsaSig.R.String(), ecdsaSig.S.String())
}

func b64urlTob64(s string) string {
	// convert base64url characters back to regular base64 alphabet
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}
	return s
}

func decode(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(b64urlTob64(s))
	if err != nil {
		panic(err)
	}
	return b
}

func b64Tob64url(s string) string {
	// convert base64url characters back to regular base64 alphabet
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.TrimRight(s, "=")
	return s
}

func encode(b []byte) string {
	return b64Tob64url(base64.StdEncoding.EncodeToString(b))
}
