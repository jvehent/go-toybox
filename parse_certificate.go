package main

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func main() {
	pemBlock, _ := pem.Decode(rawCert)
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		fmt.Println("RSA:", pub.N.BitLen(), pub.E)
	case *ecdsa.PublicKey:
		fmt.Println("ECDSA", pub.Curve.Params().Name, pub.Curve.Params().BitSize, "bits")
		fmt.Println("X:", pub.X.String())
		fmt.Println("Y:", pub.Y.String())
	}

}

var rawCert []byte = []byte(`-----BEGIN CERTIFICATE-----
MIIFQjCCAyqgAwIBAgIEAQAABDANBgkqhkiG9w0BAQwFADBqMQswCQYDVQQGEwJV
UzEQMA4GA1UEChMHQWxsaXpvbTEXMBUGA1UECxMOQ2xvdWQgU2VydmljZXMxMDAu
BgNVBAMTJ0FsbGl6b20gU2lnbmluZyBTZXJ2aWNlcyBJbnRlcm1lZGlhdGUgMTAe
Fw0xNjAzMTEyMTQ4MDRaFw0yMTAzMTAyMTQ4MDRaMIGjMQswCQYDVQQGEwJVUzET
MBEGA1UECBMKQ2FsaWZvcm5pYTEcMBoGA1UEChMTQWxsaXpvbSBDb3Jwb3JhdGlv
bjEXMBUGA1UECxMOQ2xvdWQgU2VydmljZXMxITAfBgNVBAMTGEFsbGl6b20gQ29u
dGVudCBTaWduZXIgMTElMCMGCSqGSIb3DQEJARYWaG9zdG1hc3RlckBtb3ppbGxh
LmNvbTB2MBAGByqGSM49AgEGBSuBBAAiA2IABAvUj5zuuJ8lPWb7SdomIon8XUDm
YXKzVb8wOWfE20b4a93vWBnQTbHscMfz/dMdvmTR7iCu4eK4sjBlOexlvs5hkTTr
t+WgtKrEcbN0RLICjjfl8koR6Hor7iIzi2paD6OCAWIwggFeMAwGA1UdEwEB/wQC
MAAwFAYDVR0lBA0wCwYJKwYBBAHrSQYBMB0GA1UdDgQWBBQOWlCh1q7xjcB4tC7n
PE5bYETqGTCB2AYDVR0jBIHQMIHNgBSusxmcW1IS4rh6L3/mFgMUD8oHNaGBrqSB
qzCBqDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1Nb3VudGFp
biBWaWV3MRwwGgYDVQQKExNBZGRvbnMgVGVzdCBTaWduaW5nMSQwIgYDVQQDExt0
ZXN0LmFkZG9ucy5zaWduaW5nLnJvb3QuY2ExMDAuBgkqhkiG9w0BCQEWIW9wc2Vj
K3N0YWdlcm9vdGFkZG9uc0Btb3ppbGxhLmNvbYIEAQAABDA+BglghkgBhvhCAQQE
MRYvaHR0cHM6Ly9zaWduaW5nLnNlcnZpY2VzLmFsbGl6b20ub3JnL2NhL2NybC5w
ZW0wDQYJKoZIhvcNAQEMBQADggIBAEj/AnznEmlpEWpN57JCMc+wl7VZ6TXGamwo
kQFIzlF/wUnS80ULdH4h45xglH6T4w/3L+sYjsXnML+7FQTOajtPfLqAcWQhXtpi
JdOIjx6E0592taSuzN2WlL4LZ5aKMbFrEONKjHsMpGMKeLWT3zvYppVLTumZ0HKW
AYLFDRPw1xqCwsDpz7l3y8zE5W5N5O36HLCqOn9fEA4NN197P1OCdgcYC2jO392r
lxVpL1KyflcKwPobkfuQuloBuUs+0xxD3Z6v9i3DPCnq5XN9GRGcGHmOe0IasZvy
k5rjiA/izjnIJoKQ8SC39Wy039O4wC/sITQpb9zG889Jva41OV/5BrfSQ5M/fRXr
NzSpTHOHea5tcCx4Gg4tbyW4X7zQ+OwqN6Cym+l60uYBVW0ixazYmTTecWOFRJy1
VNwo/gkpaAjAB4tcA1PmN3a6SpeCMdgb+cy5LdNhlFd6TgmHJBiU8mxOvaHeSJ6T
6m8mNjAxew5tac2i9MRtJNZ0WgKzS8J7HO5s8gUTgSj/E3NrYEbC7wkT8AY8Qkmw
pEOzxEfrezH/Xj7gcGv/+0dOz5HoZLQpHLDoAhJN36Y1WYONpGxgSMlpDoK9wd68
T8AQmAk7oBv3cjoR/3nbDVlTiGkKvU27MSuEoKGifLu00A9L40FnavmbpF3R7PTf
IaNeQuyz
-----END CERTIFICATE-----`)
