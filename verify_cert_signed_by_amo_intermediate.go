package main

import (
	"crypto/x509"
	"encoding/pem"
	"log"
)

func main() {
	ee, _ := pem.Decode([]byte(`
-----BEGIN CERTIFICATE-----
MIIFkTCCA3mgAwIBAgIIFLaErK/d9TQwDQYJKoZIhvcNAQEMBQAwZTEwMC4GA1UE
AxMnc3RhZ2luZy1zaWduaW5nLWNhLTEuYWRkb25zLm1vemlsbGEub3JnMTEwLwYJ
KoZIhvcNAQkBFiJjbG91ZC1vcHMrYWRkb25zaWduaW5nQG1vemlsbGEuY29tMB4X
DTE3MDQxODE0MzY0M1oXDTI1MDIwNzE1Mjg1MVowgYYxCzAJBgNVBAYTAlVTMQsw
CQYDVQQIEwJDQTEWMBQGA1UEBxMNTW91bnRhaW4gVmlldzEPMA0GA1UEChMGQWRk
b25zMRswGQYDVQQLExJNb3ppbGxhIEV4dGVuc2lvbnMxJDAiBgNVBAMMG3Rlc3Qt
YWRkb25AdGVzdC5tb3ppbGxhLm9yZzCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBAKF44LA+bFsPmgUmK9bmkkjJ2SfwfGkXgadTyhQ5BOjHIZXR4Cm1JROT
FggSjZNQ96JpLKSJPBU4mM5y3hx/5ygQyc3WnMKfOouPrUDXERUlsX2lJWK64RYQ
ZI2m0Q8+oZkePgQA3nDM6Ib0UGbZRRR09Ua4fU6vTlvGvrzjMyv3A8MqAG9Le0l7
mzn2Z6pU+joyKB84+sXp9EMLQL155+0UEXpKKQX0iBSM7yBX6/R6EVZjHFmNe6cy
WxZrIA5l/bn2QeAZFD8PFN03red1Jv/Ue1FFn9vTnoCNs9AcuR434V8V3i43x9Bb
y71Rt4sQ4ePLzix+fe8ATjeTXpX8y3DrpPEYlJuYZ03PTJGhzHAApDQOIRPYN8wM
n44EuNC2dHtLhbZ/UtQLgBcWYpO2Jn8QxlSoBeeLYQuU2l8woyNLBMZC0FmFp5Wz
mYfYIXFIVTKtFBMQqOSmsLMMJR4XjVkimGady6FHPHBo/TOPpnl5xSnxGGUWa2HE
p8YNQbgW1/EtvkUYvr3+PsdUM84wOuMeQ2ToyVic3zOG4Qb1Iip7810c20sz2RfQ
OPv7LmQF0HODBJd1GRI6swb5dKWoEG15vUESzZ5RSgV4Ahdj/jSYdD0Gl/SCmapy
dp8keqoJVC9tt6SisNyHfZD7oo5LuiW0dWyTI9qTfXC7WcLenlD9AgMBAAGjIzAh
MB8GA1UdIwQYMBaAFJdY6TxqKclAJNRVC7WhoqcFT3rwMA0GCSqGSIb3DQEBDAUA
A4ICAQBLMfzVKnbhKxt5rrD1Hjj29asC5tPPjwzXHziKo4FlhAiOFmIM3ztocRFN
UsC8J2E5NEr1iP9cwbbpGpaAFFuIK3xySKQs16snD9oUHyzElA0RLMXV0bjVTc3E
Dr/9GPs7Pn7j42ySwpCc6np4smlN2C/TsMMXuIYqVrFvHFSlTGMcsp3EK/jr++rL
ulmPE7RR0r/IiSS+Oqua5xLj6RAdiqmCqUZEb3KCnrNBt+kD7n/iHMxUsQHUxkd8
SH4216/IM+d3WDKWbZ7wD9rFxvP68kHYi+orcC6USXj4p1DUeyPsYr/xV2PJ0v3z
3ZpHFjjrAVCkjS/FQ1dNjw1JPSqfVsD4HmiGXknrgvlq2cAwB7yY737tftC7j6+F
IK1t7a5guTTuB2aA7+uMG4uM9tGykZ/qq09ffcNLlY0SeqqbxeA2kNTmpBeGX0VF
6LoGuI3mSD+B3EmYt7lsRCRytmVWliLjXP4HIFttM58Nmbp10pdUFJAFCX6v9klc
N5DFokr35QSiAEEYBsn9CDqrZEtsN4Tl5vK/xMvwmXmWlVkqjYRIYgtoahIQK0FA
uD3bAEfE3U+A75Ds3ymuKWxfFvGXmdFeezS8ByMhil+7SLjOdY29VBMdzP349nNx
RRxl7OEU6rdPW9iz08KLUP/uIaalOY8i5PwMSX75owtmdaJ2DQ==
-----END CERTIFICATE-----
`))
	if ee == nil {
		panic("failed to parse certificate PEM")
	}
	eecert, err := x509.ParseCertificate(ee.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	inter, _ := pem.Decode([]byte(`
-----BEGIN CERTIFICATE-----
MIIG9DCCBNygAwIBAgIEAQAAADANBgkqhkiG9w0BAQwFADCBqDELMAkGA1UEBhMC
VVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRwwGgYDVQQK
ExNBZGRvbnMgVGVzdCBTaWduaW5nMSQwIgYDVQQDExt0ZXN0LmFkZG9ucy5zaWdu
aW5nLnJvb3QuY2ExMDAuBgkqhkiG9w0BCQEWIW9wc2VjK3N0YWdlcm9vdGFkZG9u
c0Btb3ppbGxhLmNvbTAeFw0xNzA0MTQxOTU0NDZaFw0yNTAyMDcxNTI4NTFaMGUx
MDAuBgNVBAMTJ3N0YWdpbmctc2lnbmluZy1jYS0xLmFkZG9ucy5tb3ppbGxhLm9y
ZzExMC8GCSqGSIb3DQEJARYiY2xvdWQtb3BzK2FkZG9uc2lnbmluZ0Btb3ppbGxh
LmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAMsI1iJLuGfoiVNg
m07ax47Ji/JS/P2M73TZdP1QI5s4Mh9St54AA1neuc7q+lIn2XEs3vNzUtJnpFJt
UMl8lZmWII5LvcTWL73r0o3SyaBvZlrj4BuFhp0He87gsP/GKXnVs/35QLXy90t2
GffTqnnPr/qYxkOryJp2O9EE0UGOxesLzW/HA3mV30G78Zk1F3I48NpcTO/rJXm+
X+bErho+EVrNoyVBajHZi4XFyyUM3Rzszdov3//8EP1UvwuMIR8KXl7BbEFhqr82
vzm8i9X4oGassTC1OrP5Gf+B/aVGiFlQ7Dm125YENwq9pVGtupgK0UnpzAbhM68E
EO4RbxrgKZGpwUVT7GflP+/sQa4uqHm1D+DfKWW/oKeuryFtuoT/403oUkoaK94X
O/VBWzt+0zwrsYlquU97jTO5zR8aKU0ivbehJZWOrgtkCWI6U/Wsb5JViTY6q9zQ
QT6vcIr4v0E+zqabNFSYy69fTp0AFhmolF2RNEst+C484w/x1IB/RWDqW04vhg9/
cMyLcwS1NpW93Z/ws/kz3bYGxQNEXMSw6piQOHwqtEYn2Fh+jraPz4G0ASXl582e
ZKzqAcc0TyMGwIp3yBwmRhHZ/AzHGYYL67eV1jweZ+k+rJA02qERGoLXw5iMWMwO
vzVcKr72m8Q9fK/wyLLiMDEHDjTDAgMBAAGjggFmMIIBYjAMBgNVHRMEBTADAQH/
MA4GA1UdDwEB/wQEAwIBBjAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAzAdBgNVHQ4E
FgQUl1jpPGopyUAk1FULtaGipwVPevAwgdUGA1UdIwSBzTCByoAUhOpfzW9Mki/q
0fT1yKOywgy6pvahga6kgaswgagxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEW
MBQGA1UEBxMNTW91bnRhaW4gVmlldzEcMBoGA1UEChMTQWRkb25zIFRlc3QgU2ln
bmluZzEkMCIGA1UEAxMbdGVzdC5hZGRvbnMuc2lnbmluZy5yb290LmNhMTAwLgYJ
KoZIhvcNAQkBFiFvcHNlYytzdGFnZXJvb3RhZGRvbnNAbW96aWxsYS5jb22CAQEw
MwYJYIZIAYb4QgEEBCYWJGh0dHA6Ly9hZGRvbnMubW96aWxsYS5vcmcvY2EvY3Js
LnBlbTANBgkqhkiG9w0BAQwFAAOCAgEALHBbBibr/RMK3zIz+GuLOJMqljcMJNnC
+/01q1FrmHP5UIQBaPSAku0iY2m0XIpgRkvILu+4WSHEUN0stEjLCWN4Mzm3susJ
GVpi0QcFuLuVBTY/cqDEsn36GyKjmxTOH7Ke/5asgsl6XSHXoNO8HwodyqXEKPe6
wn1KQINMFuVLwVS7LZruSIF+fZOzLeWE3GsCygmXfiDTYYb4i/CIZ/coRrH7pKLX
V+kcVDSpBuQuccFeoxtvYIpSjpsV9WJJFo0QuyYDXJATeClUI6Ij5wWylzRsHUum
EnRSCjOt/XE3DJJApYFBJCQeWQU74HpFqqWZWe0lpsZFDoUUbTZDV+pMH0bYxLXr
ph6jgjMMQLOOUIybHNfbPSVmnF8t0W1/G2t0HFTpyCC+o6IFx/5KNy+Q6yT96PPH
UL3FzIMp7oq4jhceRqWMrpF3Wq3fRJTjzCkrgw+PaI4ikP0J0+UNo3KgrL1ndYl0
lTFdkpxX+KQdBfhPCglAPskcX/bb93UVnVPHf0lcHVl1yk3+YuohSDJek26HAVQd
fojAF2I6OWbDdLJQvzfHH0v5YHvGKi+GGi5Lo+mbxPEWBwv/NoA77fwka75Nkx0F
cdEbN+D+UAW6GTyNypMPCacGpA74clPKLpEWG/wpIrQne/Aap2fDX9f9hmOQs/Pa
FLcbcx7OT3I=
-----END CERTIFICATE-----
`))
	if inter == nil {
		panic("failed to parse certificate PEM")
	}
	intercert, err := x509.ParseCertificate(inter.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	root, _ := pem.Decode([]byte(`
-----BEGIN CERTIFICATE-----
MIIHYzCCBUugAwIBAgIBATANBgkqhkiG9w0BAQwFADCBqDELMAkGA1UEBhMCVVMx
CzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRwwGgYDVQQKExNB
ZGRvbnMgVGVzdCBTaWduaW5nMSQwIgYDVQQDExt0ZXN0LmFkZG9ucy5zaWduaW5n
LnJvb3QuY2ExMDAuBgkqhkiG9w0BCQEWIW9wc2VjK3N0YWdlcm9vdGFkZG9uc0Bt
b3ppbGxhLmNvbTAeFw0xNTAyMTAxNTI4NTFaFw0yNTAyMDcxNTI4NTFaMIGoMQsw
CQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExFjAUBgNVBAcTDU1vdW50YWluIFZpZXcx
HDAaBgNVBAoTE0FkZG9ucyBUZXN0IFNpZ25pbmcxJDAiBgNVBAMTG3Rlc3QuYWRk
b25zLnNpZ25pbmcucm9vdC5jYTEwMC4GCSqGSIb3DQEJARYhb3BzZWMrc3RhZ2Vy
b290YWRkb25zQG1vemlsbGEuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEAv/OSHh5uUMMKKuBh83kikuJ+BW4fQCHVZvADZh2qHNH8pSaME/YqMItP
5XQ1N5oLq1tRQO77AKn+eYPDAQkg+9VV+ct4u76YctcU/gvjieGKQ0fvuDH18QLD
hqa4DHgDmpCa/w+Eqzd54HaFj7ew9Bb7GZPHuZfk7Ct9fcN6kHneEj3KeuLiqzSV
VCRFV9RTlrUdsc1/VwF4A97JTXc3HJeWJO3azOlFpaJ8QHhmgXLLmB59HPeZ10Sf
9QwVGaKcn7yLuwtIA+wDhs8iwGZWcgmknW4DkkRDbQo7L+//4kVK+Yqq0HamZArm
vE4xENvbwOze4XYkCO3PwgmCotU7K5D3sMUUxkOaodlemO9OqRW8vJOJH3b6mhST
aunQR9/GOJ7sl4egrn2fOVZhBvM29lyBCKBffeQgtIMcKpeEKa4TNx4nTrWu1J9k
jHlvNeVL3FzMzJXRPl0RV71cYak+G6GnQ4fg3+4ZSSPxTvbwRJAO2xajkURxFSZo
sXcjYG8iPTSrDazj4LN2+882t4Q2/rMYpkowwLGbvJqHiw2tg9/hpLn1K4W18vcC
vFgzNRrTdKaJ/KjD17eJl8s8oPA7TiophPeezy1WzAc4mdlXS6A85b0mKDDU2A/4
3YmltjsSmizR2LnfeNs125EsCWxSUrAsnUYRO+lJOyNr7GGKGscCAwZVN6OCAZQw
ggGQMAwGA1UdEwQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMBYGA1UdJQEB/wQMMAoG
CCsGAQUFBwMDMCwGCWCGSAGG+EIBDQQfFh1PcGVuU1NMIEdlbmVyYXRlZCBDZXJ0
aWZpY2F0ZTAzBglghkgBhvhCAQQEJhYkaHR0cDovL2FkZG9ucy5tb3ppbGxhLm9y
Zy9jYS9jcmwucGVtMB0GA1UdDgQWBBSE6l/Nb0ySL+rR9PXIo7LCDLqm9jCB1QYD
VR0jBIHNMIHKgBSE6l/Nb0ySL+rR9PXIo7LCDLqm9qGBrqSBqzCBqDELMAkGA1UE
BhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRwwGgYD
VQQKExNBZGRvbnMgVGVzdCBTaWduaW5nMSQwIgYDVQQDExt0ZXN0LmFkZG9ucy5z
aWduaW5nLnJvb3QuY2ExMDAuBgkqhkiG9w0BCQEWIW9wc2VjK3N0YWdlcm9vdGFk
ZG9uc0Btb3ppbGxhLmNvbYIBATANBgkqhkiG9w0BAQwFAAOCAgEAck21RaAcTzbT
vmqqcCezBd5Gej6jV53HItXfF06tLLzAxKIU1loLH/330xDdOGyiJdvUATDVn8q6
5v4Kae2awON6ytWZp9b0sRdtlLsRo8EWOoRszCqiMWdl1gnGMaV7e2ycz/tR+PoK
GxHCh8rbOtG0eiVJIyRijLDjtExW8Eg+uz6Zkg1IWXqInj7Gqr23FOqD76uAfE82
YTWW3lzxpP3gL7pmV5G7ob/tIyAfrPEB4w0Nt2HEl9h7NDtKPMprrOLPkrI9eAVU
QeeI3RpAKnXOFQkqPYPXIlAaJ6qxtYa6tWHOqRyS1xKnvy/uWjEtU3tYJ5eUL1+2
vzNTdakJgkZDRdDNg0V3NYwza6BwL80VPSfqc1H6R8CU1uj+kjTlCEsoTPLeW7k5
t+lKHFMj0HZLNymgDD5f9UpI7yiOAIF0z4WKAMv/f12vnAPwmOPuOikRNOv0nNuL
RIpKO53Cd7aV5PdB0pNSPNjc6V+5IPrepALNQhKIpzoHA4oG+LlVVy4R3csPcj4e
zQQ9gt3NC2OXF4hveHfKZdCnb+BBl4S71QMYYCCTe+EDCsIGuyXWD/K2hfLD8TPW
thPX5WNsS8bwno2ccqncVLQ4PZxOIB83DFBFmAvTuBiAYWq874rneTXqInHyeCq+
819l9s72pDsFaGevmm0Us9bYuufTS5U=
-----END CERTIFICATE-----
`))
	if root == nil {
		panic("failed to parse certificate PEM")
	}
	rootcert, err := x509.ParseCertificate(root.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	if eecert.CheckSignatureFrom(intercert) != nil {
		log.Fatal("failed to verify signature")
	}
	if intercert.CheckSignatureFrom(rootcert) != nil {
		log.Fatal("failed to verify signature")
	}
	log.Println("Signature verification passes")
}