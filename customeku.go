// +build gencsr
//
// This is a simple script to generate test CSRs.
// It should not be part of the main build.
//
// PLEASE DO NOT CHECK IN
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"io/ioutil"
)

var oidExtensionRequest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 14}
var oidPoison = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 11129, 2, 4, 3}
var oidTLSFeature = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 1, 24}
var mustStapleFeatureValue = []byte{0x30, 0x03, 0x02, 0x01, 0x05}

func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("could not generate key")
		return
	}

	// Test case: duplicate must-staple
	template := &x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA1WithRSA,
		Subject: pkix.Name{
			CommonName: "not-example.com",
		},
		Attributes: []pkix.AttributeTypeAndValueSET{
			pkix.AttributeTypeAndValueSET{
				Type: oidExtensionRequest,
				Value: [][]pkix.AttributeTypeAndValue{
					[]pkix.AttributeTypeAndValue{
						pkix.AttributeTypeAndValue{
							Type:  oidTLSFeature,
							Value: []byte{0x30, 0x03, 0x02, 0x01, 0x05},
						},
					},
					[]pkix.AttributeTypeAndValue{
						pkix.AttributeTypeAndValue{
							Type:  oidTLSFeature,
							Value: []byte{0x30, 0x03, 0x02, 0x01, 0x05},
						},
					},
				},
			},
		},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, priv)
	if err != nil {
		fmt.Println("could not generate CSR")
		return
	}
	ioutil.WriteFile("duplicate_must_staple.der.csr", csrDER, 0666)
}
