package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	target := "some.target.example.net:443"
	proxy := "some.relay.example.com:3128"
	conn, err := net.DialTimeout("tcp", proxy, 30*time.Second)
	if err != nil {
		panic(err)
	}
	// send a CONNECT request to the proxy
	fmt.Fprintf(conn, "CONNECT "+target+" HTTP/1.1\r\nHost: "+target+"\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Printf("'%s'\n", status)
	// 9th character in response should be "2"
	// HTTP/1.0 200 Connection established
	//          ^
	if status == "" || len(status) < 12 {
		panic("invalid response")
	}
	if status[9] != '2' {
		fmt.Println(status)
		panic("invalid status")
	}
	fmt.Println("success")

	// do tcp stuff with `conn`

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
