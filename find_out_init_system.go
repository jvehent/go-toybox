// Find out what init system your linux is running
//
// $ go run find_out_init_system.go
// Systemd
//
// ulfr - 2014

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	initCmd, err := ioutil.ReadFile("/proc/1/cmdline")
	if err != nil {
		panic(err)
	}
	init := fmt.Sprintf("%s", initCmd)
	if strings.Contains(init, "init [") {
		fmt.Println("System-V")
	} else if strings.Contains(init, "/sbin/init") {
		fmt.Println("Upstart")
	} else if strings.Contains(init, "/usr/lib/systemd/systemd") {
		fmt.Println("Systemd")
	} else {
		fmt.Println("Can't figure out the init system")
	}
}
