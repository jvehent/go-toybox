package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IP uint32

// ToString converts an IP to its string representation
func (ip IP) ToString() (ipstr string) {
	var octets [4]string
	for i := 0; i <= 3; i++ {
		// grab one octet
		octet := ip & 255
		// move ip 8 bits to the right
		ip = ip >> 8
		octets[i] = fmt.Sprintf("%d", octet)
	}
	ipstr = octets[3] + "." + octets[2] + "." + octets[1] + "." + octets[0]
	return
}

// FromString loads an IP string into its uint32 form
func FromString(ipstr string) (ip IP, err error) {
	octets := strings.Split(ipstr, ".")
	if len(octets) != 4 {
		fmt.Errorf("Invalid IP string format")
	}
	for i := 0; i <= 3; i++ {
		octet, err := strconv.ParseUint(octets[i], 10, 8)
		if err != nil {
			return ip, err
		}
		if i > 0 {
			ip = ip << 8
		}
		ip += IP(octet)
	}
	return
}

func main() {
	ipstr := "192.168.1.2"

	fmt.Printf("Convert ip address '%s'\n", ipstr)
	ip, err := FromString(ipstr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("IP in uint32 form='%d'\n", ip)

	fmt.Printf("Back in string form='%s'\n", ip.ToString())
}
