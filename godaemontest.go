package main

import (
	"github.com/VividCortex/godaemon"
	"os"
)

func main() {
	godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	os.Exit(0)
}
