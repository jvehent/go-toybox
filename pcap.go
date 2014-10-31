// Requires:
//	# yum install libpcap-devel
// then
//	$ go get -u code.google.com/p/gopacket/pcap
//
package main

import (
	"fmt"

	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
)

func main() {
	if handle, err := pcap.OpenLive("eth0", 1600, true, 0); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter("tcp and port 80"); err != nil { // optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			fmt.Println(packet.String())
		}
	}
}
