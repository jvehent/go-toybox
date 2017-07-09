package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	//Define the AMQP binding
	bindQueue := "mig.agt.linux.usedtobeaspy1337"
	bindKey := "mig.agt.heartbeats"

	// create an AMQP configuration with specific timers
	var dialConfig amqp.Config
	dialConfig.Heartbeat = 30 * time.Second
	dialConfig.Dial = func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, 5*time.Second)
	}

	// import the client certificates
	cert, err := tls.X509KeyPair(AGENTCERT, AGENTKEY)
	if err != nil {
		panic(err)
	}

	// import the ca cert
	ca := x509.NewCertPool()
	if ok := ca.AppendCertsFromPEM(CACERT); !ok {
		panic("failed to import CA Certificate")
	}
	TLSconfig := tls.Config{Certificates: []tls.Certificate{cert},
		RootCAs:            ca,
		InsecureSkipVerify: false,
		Rand:               rand.Reader}

	dialConfig.TLSClientConfig = &TLSconfig
	// Open AMQP connection
	conn, err := amqp.DialConfig("amqps://agentspy:usedtobeaspy@publicrelay.mig.mozilla.org:443/mig", dialConfig)
	if err != nil {
		panic(err)
	}

	mqchan, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	// Limit the number of message the channel will receive at once
	err = mqchan.Qos(1, // prefetch count (in # of msg)
		0,     // prefetch size (in bytes)
		false) // is global

	_, err = mqchan.QueueDeclare(bindQueue, // Queue name
		false, // is durable
		true,  // is autoDelete
		false, // is exclusive
		false, // is noWait
		nil)   // AMQP args
	if err != nil {
		panic(err)
	}

	err = mqchan.QueueBind(bindQueue, // Queue name
		bindKey, // Routing key name
		"mig",   // Exchange name
		false,   // is noWait
		nil)     // AMQP args
	if err != nil {
		panic(err)
	}

	// Consume AMQP message into channel
	conchan, err := mqchan.Consume(bindQueue, // queue name
		"",    // some tag
		false, // is autoAck
		false, // is exclusive
		false, // is noLocal
		false, // is noWait
		nil)   // AMQP args
	if err != nil {
		panic(err)
	}
	for m := range conchan {
		// Ack this message only
		m.Ack(true)
		fmt.Printf("%s\n", m.Body)
	}

}

// CA cert that signs the rabbitmq server certificate, for verification
// of the chain of trust. If rabbitmq uses a self-signed cert, add this
// cert below
var CACERT = []byte(`-----BEGIN CERTIFICATE-----
MIIGWTCCBEGgAwIBAgIJANo/KFjOCu6fMA0GCSqGSIb3DQEBBQUAMIHRMQswCQYD
VQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNTW91bnRhaW4g
VmlldzEcMBoGA1UEChMTTW96aWxsYSBDb3Jwb3JhdGlvbjE2MDQGA1UECxMtTW96
aWxsYSBDb3Jwb3JhdGlvbiBSb290IENlcnRpZmljYXRlIFNlcnZpY2VzMRgwFgYD
VQQDEw9Nb3ppbGxhIFJvb3QgQ0ExJTAjBgkqhkiG9w0BCQEWFmhvc3RtYXN0ZXJA
bW96aWxsYS5jb20wHhcNMTIwODA0MDAwMTMyWhcNMjIwODAyMDAwMTMyWjCB0TEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDU1vdW50
YWluIFZpZXcxHDAaBgNVBAoTE01vemlsbGEgQ29ycG9yYXRpb24xNjA0BgNVBAsT
LU1vemlsbGEgQ29ycG9yYXRpb24gUm9vdCBDZXJ0aWZpY2F0ZSBTZXJ2aWNlczEY
MBYGA1UEAxMPTW96aWxsYSBSb290IENBMSUwIwYJKoZIhvcNAQkBFhZob3N0bWFz
dGVyQG1vemlsbGEuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA
uM/Xy4o0CG1S7bJrYseVHz2Dy6AjNjgFhV5p/XsVG8oddARUq7QQW2jxkaqoUanO
4AvOPH1FSkSSjTI/wh3GHSa046kZkqDBdjF6MV+MgS2aaEphbv1jU4NfY7XNMamo
yzh5uaI4z4ktSx0k/m9y+GlCs+R4jE/U3GG8TQRBz7cg+mA3XWAEUlaXE6bGDJIg
YjKpDz3MSwJgu3csOwimQOelaERF6gL0bamsSXRVEmSJnLwU4hZ0q/pDg4pf7D4x
rVdKey8V3PHiCMotAVayWISzhatsgRX8Kh6t1GvsChArm/B5VuQNcxYqJh1VJQgM
IrOccmdZKK2M3yc56K8RzifejB97SU7vLsgXESKe6COYEoHgaV3vpsdHp2pjeuPm
CgDqXLzEHB2LnVSq7O5DTM7HnV6+k26KpI2V0S6IjlduYjrtc2QcyG3oe+TkD7AP
+PYRdY+OoNhj/AmgnIdktNyjLlQ0m9c6+zGoUzRVeLcA1rrcgB8wx2v24JeS0R/y
yEZ6IBAIRYggl1YEYxM66JneGrrZd2K2YxO8IapOOoD8NZYh+ltESweKofK0ZDSX
bHIsKTom/hmXflkH6shwvS4iZHucib93KIRZRWu/sceGsGBntWJaj8+JaaeVV1zv
Plq7JT3aIEeZs4iYSL69lUOENo2k10diYng2cBIG8UUCAwEAAaMyMDAwDwYDVR0T
AQH/BAUwAwEB/zAdBgNVHQ4EFgQUudVz7I8Z4rMZOuRuoYZQNYtmwukwDQYJKoZI
hvcNAQEFBQADggIBABcQo7pwie4XlfvPAPteBTTiMFfhhETEnBixkBVRkXCdDgXz
HMAeswPATy3JMM7oKG+1fb7dxYdyCMMdPwIAg5npxKJbYbKSkhU94O4Hj4DtUIHY
AkCqBAu0GvUCBsRehbdMm32e1A6Py+Y9SixjPFq6zOKzLO8p+YAPFymhNP345yp9
STUDffKMqQIUkG51LusjQkysAgFtHAO2s1p9gBGgtlXmBS2xcrBJvXUBMHshYAo+
igSBQe8d2dQHr1MMop7i4bkK8mKFdLBlmpJ1NF7X+RRGLr7kCRTt6kbjGUkTHIVf
js5Q14oHwU8P5A0yKZusBY9b94ZeiJhPKD9NzsDYj838EVs0h3FfPnS+ymePsKDO
B66FgITGE/j7Z3LxroYUgVb6RHrBES/KYOkH7UqVwtjWFqKq8UPNeH2wD78vBVEz
NJuiujLSxnnkW0rC6GuOAipGl84BrX5/IPglD+NkqzpqdiBsRQ/svMuHiHtDKma/
7p0qWlE3q9FmlbG3lZ25o92fUQCMIiuEwOQEhj9PAGUyv/nduGNftu2Ma3L8bGFl
XORWVs12ww/M5WUSnft/3Lpi05ZTDnoXndrZ3TkHo3jNzqPGZ8ifE6LBM5MuAMmd
8BYb9xs9ANv6R8ivcX5Phcycda0Kez6p0iUcZBDibMWQqPp83YfhuAgtkhnM
-----END CERTIFICATE-----`)

// All clients share a single X509 certificate, for TLS auth on the
// rabbitmq server. Add the public client cert below.
var AGENTCERT = []byte(`-----BEGIN CERTIFICATE-----
MIIFEjCCAvqgAwIBAgICAp0wDQYJKoZIhvcNAQEFBQAwgdExCzAJBgNVBAYTAlVT
MRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRww
GgYDVQQKExNNb3ppbGxhIENvcnBvcmF0aW9uMTYwNAYDVQQLEy1Nb3ppbGxhIENv
cnBvcmF0aW9uIFJvb3QgQ2VydGlmaWNhdGUgU2VydmljZXMxGDAWBgNVBAMTD01v
emlsbGEgUm9vdCBDQTElMCMGCSqGSIb3DQEJARYWaG9zdG1hc3RlckBtb3ppbGxh
LmNvbTAeFw0xNTA3MjAxODQ0MTBaFw0yNTA3MTcxODQ0MTBaMFIxKjAoBgNVBAMU
IWFnZW50LW9wc2VjQG1pZy5vcHNlYy5tb3ppbGxhLm9yZzEkMCIGCSqGSIb3DQEJ
ARYVb3BzZWMrbWlnQG1vemlsbGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAuw+dYYs2gK0InJD9rS3FG19ESnRZH+tU2pqg6TLM4K2krfxy7PGC
hDL8/kPzTTDYHzCRgaqBXTC5qphcoyM7QhduyWXjjt85FCTyL3/zQhZeuRO3d0WX
FaL+lB/yX6G9CSs/BhLJ7f3jSoesk5VUjCwGYkq4P2XEVCZ42Xd3YEx3Vk1I5Pwa
jEPNeIIkiLRwbtGmdNg13FR832Z4tjDzaRXrMU3db0fQBLFwtwkA+2bZiwEKtyLm
ZX0FQTQKxbXRs7+UCEzl4OT3ja/4o5z7T0snuDnbUGuZRodLJc3QZWNCzwxLoI5f
PKCnkqjvk2zxrSXKEa5Q7oWFT+vR9aTwawIDAQABo3IwcDAgBgNVHREEGTAXgRVv
cHNlYyttaWdAbW96aWxsYS5jb20wDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBS5
1XPsjxnisxk65G6hhlA1i2bC6TAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUH
AwQwDQYJKoZIhvcNAQEFBQADggIBAFY7Ze6y15av6u80sa2KLumWFYg5E7KoT9gb
s+n1WNtDVkLkndcVh0WGcZW8YbQDdVQh255TBMHq67W+MXXqLle5Ws8yYSNJCmyc
wvkKRzGCroaXN3A+tqqHHBF1ckUs7FqSdm/HRNkvgCc7EEl/Bzchuhvml4LOiGu1
mQWlGKDe6mxvykkeWA8emnxbBQddiL5KTFdZRVQa6Il8GQb5gL2MkLsn2JyS6mKk
RPbupDMfSfz4TDCM+zqbmJOrJvqhuoeIQT0586uQ/eZA9sAd7w0mGgJgzcRyN6Ur
7eE25DsAxMhlp8FuT+KK5dd6mSpdRl5RHU6YnEYomZpqQN+k0/3yRwsiKUygGq0q
s54x4nPSZxRwgdJCQ/jFIugvA/hmmKaVQ0bDSV+rE4HCp/uEtTuHnb06i/KZs33z
+2nzjycWVqqMRqEUuHzGSxy1LEA5KgbB7ifDz9Uhy9B71/g+17ainx0u5CJ4PmJL
2D48fssIzZ3dUTGeRggPz8E/dg6fl4rX3STo8yrNQm1za/04uBwHFdqSa0ixT7GP
yvsRzdm1+f97NxT6H4d7xBy1gi0ziRV8QzI5Z3PzLcilBKsn5MTkOV93Y5FVq+K2
2JpFSiJScQqr8JbyYV8sfKQHrfFM2+N33oUH3xTU5LyhzTeQEUBlhaKWDOQJan0G
R1Qc5uXF
-----END CERTIFICATE-----`)

// Add the private client key below.
var AGENTKEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAuw+dYYs2gK0InJD9rS3FG19ESnRZH+tU2pqg6TLM4K2krfxy
-----END RSA PRIVATE KEY-----`)
