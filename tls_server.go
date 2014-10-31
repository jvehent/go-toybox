package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		panic(err)
	}
	config := tls.Config{
		Certificates:             []tls.Certificate{certificate},
		ClientAuth:               tls.RequireAnyClientCert,
		MinVersion:               tls.VersionTLS10,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA},
	}
	config.Rand = rand.Reader

	netlistener, err := tls.Listen("tcp", "127.0.0.1:50443", &config)
	if err != nil {
		panic(err)
	}
	newnetlistener := tls.NewListener(netlistener, &config)
	fmt.Println("I am listening...")
	for {
		newconn, err := newnetlistener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Got a new connection from %s. Say Hi!\n", newconn.RemoteAddr())
		newconn.Write([]byte("ohai"))
		newconn.Close()
	}
}
