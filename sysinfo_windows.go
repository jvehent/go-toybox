package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var hostname, domain, osname, osversion string
	fmt.Println("get data from the systeminfo")
	out, err := exec.Command("systeminfo").Output()
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(out)
	reader := bufio.NewReader(buf)
	iter := 0
	fmt.Println("parsing systeminfo output")
	for {
		lineBytes, _, err := reader.ReadLine()
		// don't loop more than 2000 times
		if err != nil || iter > 2000 {
			goto exit
		}
		line := fmt.Sprintf("%s", lineBytes)
		fmt.Println(line)
		if strings.Contains(line, "OS Name:") {
			out := strings.SplitN(line, ":", 2)
			if len(out) == 2 {
				osname = out[1]
			}
			osname = cleanString(osname)
		} else if strings.Contains(line, "OS Version:") {
			out := strings.SplitN(line, ":", 2)
			if len(out) == 2 {
				osversion = out[1]
			}
			osversion = cleanString(osversion)
		} else if strings.Contains(line, "Domain:") {
			out := strings.SplitN(line, ":", 2)
			if len(out) == 2 {
				domain = out[1]
			}
			domain = cleanString(domain)
		} else if strings.Contains(line, "Host Name:") {
			out := strings.SplitN(line, ":", 2)
			if len(out) == 2 {
				hostname = out[1]
			}
			hostname = cleanString(hostname)
		}
		iter++
	}
exit:
	fmt.Printf("hostname: '%s'\ndomain: '%s'\nOS name: '%s'\nOS Version: '%s'\n",
		hostname, domain, osname, osversion)
	time.Sleep(10 * time.Second)
}

// cleanString removes spaces, quotes and newlines
func cleanString(str string) string {
	if len(str) < 1 {
		return str
	}
	if str[len(str)-1] == '\n' {
		str = str[0 : len(str)-1]
	}
	// remove heading whitespaces and quotes
	for {
		if len(str) < 2 {
			break
		}
		switch str[0] {
		case ' ', '"', '\'':
			str = str[1:len(str)]
		default:
			goto trailing
		}
	}
trailing:
	// remove trailing whitespaces, quotes and linebreaks
	for {
		if len(str) < 2 {
			break
		}
		switch str[len(str)-1] {
		case ' ', '"', '\'':
			str = str[0 : len(str)-1]
		default:
			goto exit
		}
	}
exit:
	return str
}
