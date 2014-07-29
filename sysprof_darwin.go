package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// get data from the system_profiler
	// don't attempt to use their weird xml format, better parse plain text
	out, err := exec.Command("system_profiler", "SPSoftwareDataType").Output()
	if err != nil {
		panic(err)
	}
	var sysv, kerv string
	buf := bytes.NewBuffer(out)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			goto exit
		}
		// chomp
		if line[len(line)-1] == '\n' {
			line = line[0 : len(line)-1]
		}
		if strings.Contains(line, "System Version") {
			sysv = strings.Split(line, ":")[1]
			if sysv[0] == ' ' {
				sysv = sysv[1:len(sysv)]
			}
		} else if strings.Contains(line, "Kernel Version") {
			kerv = strings.Split(line, ":")[1]
			if kerv[0] == ' ' {
				kerv = kerv[1:len(kerv)]
			}
		}
	}
exit:
	fmt.Printf("system version '%s' kernel version '%s'\n", sysv, kerv)
}
