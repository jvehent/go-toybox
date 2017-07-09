package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/fullsailor/pkcs7"
)

func main() {
	cert, err := createCertificate()
	if err != nil {
		log.Fatal(err)
	}
	content := []byte(`Signature-Version: 1.0
MD5-Digest-Manifest: PQX/zDhGcHy8EobBktpx0g==
SHA1-Digest-Manifest: a3kH97BClCeW08eyhLoarQ6GQjw=

`)
	toBeSigned, err := pkcs7.NewSignedData(content)
	if err != nil {
		log.Fatalf("Cannot initialize signed data: %s", err)
	}
	if err := toBeSigned.AddSigner(cert.Certificate, cert.PrivateKey, pkcs7.SignerInfoConfig{}); err != nil {
		log.Fatalf("Cannot add signer: %s", err)
	}
	detachedSig, err := toBeSigned.FinishDetached()
	if err != nil {
		log.Fatalf("Cannot finish signing data: %s", err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(detachedSig))

	err = ioutil.WriteFile("/tmp/detachedsig.rsa", detachedSig, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p7, err := pkcs7.Parse(detachedSig)
	if err != nil {
		log.Fatalf("Cannot parse our signed data: %s", err)
	}
	p7.Content = content
	if bytes.Compare(content, p7.Content) != 0 {
		log.Fatalf("Our content was not in the parsed data:\n\tExpected: %s\n\tActual: %s", content, p7.Content)
	}
	if err := p7.Verify(); err != nil {
		log.Fatalf("Cannot verify our signed data: %s", err)
	}

}

type certKeyPair struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

func createCertificate() (certKeyPair, error) {
	signer, err := createCertificateByIssuer("Eddard Stark", nil)
	if err != nil {
		return certKeyPair{}, err
	}
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: signer.Certificate.Raw})
	pair, err := createCertificateByIssuer("Jon Snow", signer)
	if err != nil {
		return certKeyPair{}, err
	}
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: pair.Certificate.Raw})
	return *pair, nil
}

func createCertificateByIssuer(name string, issuer *certKeyPair) (*certKeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 32)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:       serialNumber,
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	var issuerCert *x509.Certificate
	var issuerKey crypto.PrivateKey
	if issuer != nil {
		issuerCert = issuer.Certificate
		issuerKey = issuer.PrivateKey
	} else {
		issuerCert = &template
		issuerKey = priv
	}
	cert, err := x509.CreateCertificate(rand.Reader, &template, issuerCert, priv.Public(), issuerKey)
	if err != nil {
		return nil, err
	}
	leaf, err := x509.ParseCertificate(cert)
	if err != nil {
		return nil, err
	}
	return &certKeyPair{
		Certificate: leaf,
		PrivateKey:  priv,
	}, nil
}
